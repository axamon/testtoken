// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/axamon/token"
)

// accesso is an interface to manage credentials.
type accesso interface {
	// autenticato method returns true whether credentials
	// are found in any storage (json file or sql db).
	autenticato(context.Context) bool

	// token method returns a psuedo token if credentials are good.
	token(context.Context) string
}

// verifica function verifies that credentials are found.
// func verifica(a accesso) {
// 	fmt.Println(a.autenticato())
// }

// getToken function returns a pseudo token.
func getToken(ctx context.Context, a accesso) string {
	return a.token(ctx)
}

// autenticato returns true if credentials are found in any storage.
func (c credentials) autenticato(ctx context.Context) bool {

	ctx, cancel := context.WithTimeout(ctx, 15*time.Millisecond)
	defer cancel()
	defer runtime.GC()

	var cc token.Credentials
	cc.User = c.User
	cc.Hashpass = c.Hashpass

	// Istanzia un wait group per gestire i processi paralleli.
	var wg sync.WaitGroup

	var globallyAuthenticated bool = false

	// Aggiunge un processo parallelo.
	wg.Add(1)
	go func() {
		defer runtime.Gosched()
		defer wg.Done()
		isAuthenticated, err := token.TestSearch(ctx, &cc)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		defer log.Printf("Finito controllo su DB:\t%v,\tid:\t%s\n", isAuthenticated, ctx.Value(k))
		if isAuthenticated {
			globallyAuthenticated = true
		}
		return
	}()

	// Aggiunge un processo parallelo.
	wg.Add(1)
	go func() {
		defer runtime.Gosched()
		defer wg.Done()
		isAuthenticated, err := token.CheckLocalCredentials(ctx, &cc)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		defer log.Printf("Finito controllo su File:\t%v,\tid:\t%s\n", isAuthenticated, ctx.Value(k))
		if isAuthenticated {
			globallyAuthenticated = true
		}
		return
	}()

	// Aspetta che tutti i processi paralleli terminino.
	wg.Wait()

	select {
	case <-ctx.Done():
		log.Printf("Error in autentiato function: %v\n", ctx.Err())
		return false
	default:
		return globallyAuthenticated
	}
}

func (c credentials) token(ctx context.Context) string {

	if c.autenticato(ctx) {
		token, err := token.GenerateToken(ctx)
		if err != nil {
			log.Println(err)
		}
		return token
	}
	return ""
}

// getUUDI generates a string to use as context-value.
func getUUDI(ctx context.Context) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Microsecond)
	defer cancel()
	defer runtime.GC()

	b := make([]byte, 16)
	_, err := rand.Read(b)

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("impossible to generate UDDI: %v", ctx.Err())

	default:
		return uuid, err
	}

}

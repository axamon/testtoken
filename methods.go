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
	"testtoken/token"
)

// accesso is an interface to manage credentials.
type accesso interface {
	// autenticato method returns true whether credentials
	// are found in any storage (json file or sql db).
	autenticato() bool

	// token method returns a psuedo token if credentials are good.
	token() string
}

// verifica function verifies that credentials are found.
func verifica(a accesso) {
	fmt.Println(a.autenticato())
}

// getToken function returns a pseudo token.
func getToken(a accesso) string {
	return a.token()
}

// autenticato returns true if credentials are found in any storage.
func (c credentials) autenticato() bool {

	var cc token.Credentials
	cc.User = c.User
	cc.Hashpass = c.Hashpass

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer runtime.GC()

	type ctxINTERFACE string
	var k ctxINTERFACE

	uddi, err := getUUDI(ctx)
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.WithValue(ctx, k, uddi)

	var wg sync.WaitGroup

	var globallyAuthenticated bool = false

	wg.Add(1)

	go func() {
		defer runtime.Gosched()
		defer wg.Done()
		isAuthenticated, err := token.TestSearch(ctx, &cc)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		defer log.Println("Finito controllo su DB", isAuthenticated, ctx.Value(k))
		if isAuthenticated {
			globallyAuthenticated = true
		}
		return
	}()

	wg.Add(1)
	go func() {
		defer runtime.Gosched()
		defer wg.Done()
		isAuthenticated, err := token.CheckLocalCredentials(ctx, &cc)
		defer log.Printf("Finito controllo su File %v, id: %s\n", isAuthenticated, ctx.Value(k))
		if err != nil {
			log.Printf("Error: %v", err)
		}
		if isAuthenticated {
			globallyAuthenticated = true
		}
		return
	}()

	// waits until both go routines have finished working.
	wg.Wait()

	return globallyAuthenticated
}

func (c credentials) token() string {

	if c.autenticato() {
		token, err := token.GenerateToken(context.TODO())
		if err != nil {
			log.Println(err)
		}
		return token
	}
	return ""
}

// getUUDI generates a string to use as context-value.
func getUUDI(ctx context.Context) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Millisecond)
	defer cancel()
	defer runtime.GC()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("impossible to generate UDDI: %v", ctx.Err())

	default:
		b := make([]byte, 16)
		_, err := rand.Read(b)

		uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

		return uuid, err
	}

}

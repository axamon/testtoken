// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"context"
	"log"
	"runtime"
	"sync"

	"github.com/axamon/uddi"
)

type ctxINTERFACE string

var k ctxINTERFACE

func init() {

	ctx := context.Background()

	udditoken, err := uddi.CreateCtx(ctx)
	if err != nil {
		log.Println(err)
	}

	ctx = context.WithValue(ctx, k, udditoken)
}

// accesso is an interface to manage credentials.
type accesso interface {
	// autenticato method returns true whether credentials
	// are found in any storage (json file or sql db).
	Autenticato(context.Context) bool

	// token method returns a pseudo token if credentials are good.
	GetToken(context.Context) string
}

// Autenticato returns true if credentials are found in any storage.
func (c Credentials) Autenticato(ctx context.Context) bool {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer runtime.GC()

	// Istanzia un wait group per gestire i processi paralleli.
	var wg sync.WaitGroup

	var globallyAuthenticated bool = false

	// Aggiunge un processo parallelo.
	wg.Add(1)
	go func() {
		defer runtime.Gosched()
		defer wg.Done()
		isAuthenticated, err := CheckCredentialsDBCtx(ctx, &c)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		defer log.Printf("Finito controllo su DB:\t%v,\tid:\t%s\n",
			isAuthenticated,
			ctx.Value(k))

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
		isAuthenticated, err := CheckLocalCredentials(ctx, &c)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		defer log.Printf("Finito controllo su File:\t%v,\tid:\t%s\n",
			isAuthenticated,
			ctx.Value(k))

		if isAuthenticated {
			globallyAuthenticated = true
		}
		return
	}()

	// Aspetta che tutti i processi paralleli terminino.
	wg.Wait()

	select {
	case <-ctx.Done():
		log.Printf("Error in autenticato %v function: %v\n",
			ctx.Value(k),
			ctx.Err())
		return false
	default:
		return globallyAuthenticated
	}
}

// GetToken ...
func (c Credentials) GetToken(ctx context.Context) string {

	if c.Autenticato(ctx) {
		token, err := GenerateCtx(ctx)
		if err != nil {
			log.Println(err)
		}
		return token
	}
	return ""
}

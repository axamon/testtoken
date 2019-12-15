// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"testtoken/token"
	"testtoken/token/db"
	"time"
)

var user, pass string

func init() {
	flag.StringVar(&user, "u", "", "Username")
	flag.StringVar(&pass, "p", "", "Password")

	flag.Parse()
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer runtime.GC()

	var c = token.Credentials{User: user, Hashpass: hashpassword(pass)}

	var globallyAuthenticated bool = false

	var wg sync.WaitGroup

	type ctxINTERFACE string
	var k ctxINTERFACE

	uddi, err := getUUDI(ctx)
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.WithValue(ctx, k, uddi)

	wg.Add(1)
	go func() {
		defer runtime.Gosched()
		defer wg.Done()
		isAuthenticated, err := db.TestSearch(ctx, &c)
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
		isAuthenticated, err := token.CheckLocalCredentials(ctx, &c)
		defer log.Println("Finito controllo su File", isAuthenticated, ctx.Value(k))
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

	switch globallyAuthenticated {
	case true:
		log.Printf("%s autenticato\n", user)
		token, err := token.GenerateToken(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(token)
		os.Exit(0)

	default:
		log.Fatalf("%s non autenticato\n", user)
		// os.Exit(1) is added automatically by log.Fatalf
	}
}

func hashpassword(s string) string {
	h := md5.New()
	h.Write([]byte(s + "\n"))
	hashedpass := h.Sum(nil)

	return fmt.Sprintf("%x", hashedpass)
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

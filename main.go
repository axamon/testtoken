// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"testtoken/token"
	"testtoken/token/db"
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

	wg.Add(1)
	go func() {
		defer runtime.Gosched()
		defer wg.Done()
		isAuthenticated, err := db.TestSearch(ctx, &c)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		defer fmt.Println("Finito controllo su DB", isAuthenticated)
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
		defer fmt.Println("Finito controllo su File", isAuthenticated)
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

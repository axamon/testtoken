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

	var authenticated = make(chan bool, 2)

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
		authenticated <- isAuthenticated
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
		authenticated <- isAuthenticated
		return
	}()

	// waits until something is put in channel
	// if the first result is true it is enough to authenticate-
	// it avoids waiting for both routines in case one has already authenticated.
	select {
	case r := <-authenticated:
		if r == true {
			fmt.Println("autenticato")
			os.Exit(0)
		}
		if r == false {
			break
		}
	}

	// waits until both go routines have finished working.
	wg.Wait()

	// another select to check if any go routine has managed to authenticate.
	select {
	case r := <-authenticated:
		if r == true {
			fmt.Println("autenticato")
			os.Exit(0)
		}
		if r == false {
			break
		}
	}
	os.Exit(1)
}

func hashpassword(s string) string {
	h := md5.New()
	h.Write([]byte(s + "\n"))
	hashedpass := h.Sum(nil)

	return fmt.Sprintf("%x", hashedpass)
}

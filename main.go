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

type accesso interface {
	autenticato() bool
	token() string
}

type credentials token.Credentials

func verifica(a accesso) {
	fmt.Println(a.autenticato())
}

func getToken(a accesso) {
	fmt.Println(a.token())
}

func main() {

	// var c = token.Credentials{User: user, Hashpass: hashpassword(pass)}

	poldo := credentials{User: "Poldo", Hashpass: hashpassword("panino")}
	var dinamico = credentials{User: user, Hashpass: hashpassword(pass)}

	verifica(poldo)
	verifica(dinamico)

	getToken(poldo)
	getToken(dinamico)

}

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
		isAuthenticated, err := db.TestSearch(ctx, &cc)
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

	// switch globallyAuthenticated {
	// case true:
	// 	log.Printf("%s autenticato\n", user)
	// 	token, err := token.GenerateToken(ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(token)
	// 	os.Exit(0)

	// default:
	// 	log.Fatalf("%s non autenticato\n", user)
	// 	// os.Exit(1) is added automatically by log.Fatalf

	// }
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

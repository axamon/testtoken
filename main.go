// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testtoken/token"
	"time"
)

// user and pass variables.
var user, pass string

var t int

func init() {
	flag.StringVar(&user, "u", "", "Username")
	flag.StringVar(&pass, "p", "", "Password")
	flag.IntVar(&t, "t", 5, "Timeout in millisecons")

	// parse cli arguments into the variables.
	flag.Parse()
}

// credentials type is needed to mask in main package the
// type created in the token package.
type credentials token.Credentials

func main() {
	// Set the main context with t milliseconds timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Millisecond)
	// At the end of function cleans up.
	defer cancel()

	// Creates the channel where the token will be passed.
	var result = make(chan string, 1)

	// Creates the credential variable to test.
	var dinamico = credentials{User: user, Hashpass: hashstring.Md5Sum(pass)}
	result <- getToken(dinamico)

	select {
	// If checks take too long it quits.
	case <-ctx.Done():
		log.Fatal(ctx.Err())
	case t := <-result:
		// Prints the pseudo token.
		fmt.Println(t)
	}

}

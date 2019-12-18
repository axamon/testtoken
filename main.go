// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/axamon/hashstring"
	"github.com/axamon/token"
	"github.com/axamon/uddi"
)

// user and pass variables.
var user, pass string

// t is the timeout to use.
var t int

// parse cli arguments into the variables.
func init() {
	flag.StringVar(&user, "u", "", "Username")
	flag.StringVar(&pass, "p", "", "Password")
	flag.IntVar(&t, "t", 30, "Timeout in millisecons")

	flag.Parse()
}

// credentials type is needed to mask in main package the
// type created in the token package.
type credentials token.Credentials

type ctxINTERFACE string

var k ctxINTERFACE

func main() {
	// Set the main context with t milliseconds timeout.
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(t)*time.Millisecond)
	// At the end of function cleans up.
	defer cancel()

	udditoken, err := uddi.CreateCtx(ctx)
	if err != nil {
		log.Println(err)
	}

	ctx = context.WithValue(ctx, k, udditoken)

	// Creates the channel where the token will be passed.
	var result = make(chan string, 1)

	// Creates the credential variable to test.
	var dinamic = credentials{User: user, Hashpass: hashstring.Md5Sum(pass)}

	// gets a pseudo UDDI token if the credentials are present in any storage.
	result <- dinamic.token(ctx)

	select {
	// If checks took too long it quits.
	case <-ctx.Done():
		log.Fatalf("Error main func timeout: %v\n", ctx.Err()) // implicitly does os.Exit(1)
	case t := <-result:
		// closes the channel.
		close(result)
		// Prints the pseudo token.
		fmt.Println(t)
		os.Exit(0)
	}

}

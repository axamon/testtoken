// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"runtime"
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

	h := md5.New()
	io.WriteString(h, pass)
	hashedpass := h.Sum(nil)

	autheticated, err := db.TestSearch(ctx, user, fmt.Sprintf("%x", hashedpass))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if !autheticated {
		fmt.Println(autheticated, "NOT AUTHENTICATED")
	}

	var c = token.Credentials{User: user, Pass: pass}

	token, err := token.GetToken(ctx, &c)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(token)

}

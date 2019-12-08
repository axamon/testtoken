// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
	t "testtoken/token"
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

	var c = t.Credentials{User: user, Pass: pass}

	token, err := t.GetToken(ctx, &c)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(token)

}

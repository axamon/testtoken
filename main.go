// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"testtoken/token"
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

	poldo := credentials{User: "Poldo", Hashpass: hashpassword("panino")}
	var dinamico = credentials{User: user, Hashpass: hashpassword(pass)}

	verifica(poldo)
	verifica(dinamico)

	getToken(dinamico)

}

func hashpassword(s string) string {
	h := md5.New()
	h.Write([]byte(s + "\n"))
	hashedpass := h.Sum(nil)

	return fmt.Sprintf("%x", hashedpass)
}

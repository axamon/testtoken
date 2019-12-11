// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

// Credentials is the type used to pass username and password around.
type Credentials struct {
	User     string
	Pass     string
	Hashpass string
}

// credentialsDB maps the json db to struct.
type credentialsDB struct {
	UserpassDB []struct {
		PasswordDB string `json:"pass"`
		UsernameDB string `json:"user"`
	} `json:"credentials"`
}

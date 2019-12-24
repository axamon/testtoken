// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token generates pseudo uddi tokens if credentials used
// match any storage.
package token

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	rand "math/rand"
	"runtime"
)

// CredentialsJSONFile is the json file containing credentials.
var CredentialsJSONFile = "credentialsdb.json"

var src cryptoSource

func init() {
	rnd := rand.New(src)
	rnd.Seed(rnd.Int63())
}

// GenerateCtx generates a token.
func GenerateCtx(ctx context.Context) (string, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer runtime.GC()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("function GenerateCtx in error: %v", ctx.Err())

	default:
		b := make([]byte, 16)
		_, err := rand.Read(b)

		uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

		return uuid, err
	}

}

// CheckLocalCredentials verifies username and passwords on local json file.
func CheckLocalCredentials(ctx context.Context, c *Credentials) (bool, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	body, err := ioutil.ReadFile(CredentialsJSONFile)

	var db = new(credentialsDB)
	err = json.Unmarshal(body, &db)
	if err != nil {
		return false, fmt.Errorf(
			"Error in unmarshalling %s: %v", CredentialsJSONFile, err)
	}

	for _, r := range db.UserpassDB {
		if r.UsernameDB == c.User && r.PasswordDB == c.Hashpass {
			return true, nil
		}
	}

	select {
	case <-ctx.Done():
		return false, fmt.Errorf(
			"function checkCredentials in error: %v", ctx.Err())
	default:
		return false, nil
	}
}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatalf("Low entropy, cannot create crypto random number: %v", err)
	}
	return v
}

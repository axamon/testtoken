// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"
	"time"

	"testtoken/token"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
)

var userdb, passdb, addr, d string

func init() {
	userdb = "pippo"
	passdb = "pippo"
	addr = "127.0.0.1:3306"
	d = "app"
}

// QueryCredentials is a query that executed on DB returns true if both user
// and password are found in a record.
const QueryCredentials = "SELECT IF(COUNT(*),'true','false') FROM app.credentials WHERE username = ? AND password = ?"

func TestSearch(ctx context.Context, c *token.Credentials) (bool, error) {

	// Crea il contesto base.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("Timeout: %v", ctx.Err())
	default:

		// Apre la conessione al DB.
		db, err := sql.Open("mysql", userdb+":"+passdb+"@tcp("+addr+")/"+d)
		// Chiude la connesione al DB alla fine.
		defer db.Close()

		err = db.PingContext(ctx)
		if err != nil {
			return false, fmt.Errorf("DB not rechable: %v", err)
		}

		// Se c'è un errore esce.
		if err != nil {
			return false, fmt.Errorf("db access not possible: %v", err)
		}

		var isAuthenticated bool
		err = db.QueryRowContext(ctx, QueryCredentials, c.User, c.Hashpass).Scan(&isAuthenticated)
		if err != nil {
			return false, fmt.Errorf("db access not possible: %v", err)
		}

		return isAuthenticated, err
	}
}

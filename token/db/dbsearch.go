package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
	"testtoken/token"
)

var userdb, passdb, addr, d string

const querycredentials = "SELECT IF(COUNT(*),'true','false') FROM app.credentials where username = ? and password = ?"


func TestSearch(ctx context.Context, c *token.Credentials) (bool, error) {

	// Crea il contesto base.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userdb = "pippo"
	passdb = "pippo"
	addr = "127.0.0.1:3306"
	d = "app"
	// Apre la conessione al DB.
	db, err := sql.Open("mysql", userdb+":"+passdb+"@tcp("+addr+")/"+d)

	// Se c'è un errore esce.
	if err != nil {
		return false, fmt.Errorf("db access not possible: %v", err)
	}

	// Chiude la connesione al DB alla fine.
	defer db.Close()

	var isAuthenticated bool
	err = db.QueryRowContext(ctx, querycredentials, c.User, c.Hashpass).Scan(&isAuthenticated)
	if err != nil {
		log.Fatal(err)
	}

	return isAuthenticated, err
}

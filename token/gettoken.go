package token

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
)

const gettokeninerror = "GetToken function in error: %v"
const genertateTokenInError = "function gerateToken in error: %v"
const checkcredentialsinerror = "function checkCredentials in error: %v"

const credentialsdb = "credentialsdb.json"

// GetToken generate a uuid like token (does not follow standards).
func GetToken(ctx context.Context, c *Credentials) (s string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	defer runtime.GC()

	var errors = make(chan error, 1)

	err = checkCredentials(ctx, c)
	if err != nil {
		errors <- err
	}

	select {
	case err = <-errors:
		return "", err

	case <-ctx.Done():

		return "", fmt.Errorf(gettokeninerror, ctx.Err())

	default:
		s, err = generateToken(ctx)
		if err != nil {
			errors <- err
		}

	}

	return s, err
}

// checkCredentials verifies username and passwords.
func checkCredentials(ctx context.Context, c *Credentials) error {

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	defer runtime.GC()

	var err error
	var errors = make(chan error, 1)

	select {
	case err = <-errors:
		return err

	case <-ctx.Done():
		return fmt.Errorf(checkcredentialsinerror, ctx.Err())

	default:
		body, err := ioutil.ReadFile(credentialsdb)

		var db = new(credentialsDB)
		err = json.Unmarshal(body, &db)
		if err != nil {
			errors <- err
		}

		for _, r := range db.userpassDB {
			if r.usernameDB == c.User && r.passwordDB == c.Pass {
				return nil
			}
		}
	}

	return fmt.Errorf("bad credentials")
}

// generateToken generates a token.
func generateToken(ctx context.Context) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Millisecond)
	defer cancel()
	defer runtime.GC()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf(genertateTokenInError, ctx.Err())

	default:
		b := make([]byte, 16)
		_, err := rand.Read(b)

		uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

		return uuid, err
	}

}

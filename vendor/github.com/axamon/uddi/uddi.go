// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uddi creates pseudo UDDI tokens.
package uddi

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// CreateCtx generates a UDDI exadedimal string.
func CreateCtx(ctx context.Context) (string, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	b := make([]byte, 16)
	_, err := rand.Read(b)

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("impossible to generate UDDI: %v", ctx.Err())

	default:
		return uuid, err
	}
}

// Create generates a UDDI exadedimal string.
func Create() string {

	b := make([]byte, 16)
	rand.Read(b)

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

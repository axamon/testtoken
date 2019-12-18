// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"testing"

	"github.com/axamon/hashstring"
)

func Test_credentials_autenticato(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		c    credentials
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "primo", c: credentials{User: "pippo", Hashpass: hashstring.Md5Sum("pippo")}, args: args{ctx: context.TODO()}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.autenticato(tt.args.ctx); got != tt.want {
				t.Errorf("credentials.autenticato() = %v, want %v", got, tt.want)
			}
		})
	}
}

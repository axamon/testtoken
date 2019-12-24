// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"context"
	"github.com/axamon/hashstring"
	"github.com/axamon/token"
	"math/rand"
	"testing"
)

func TestCredentials_token(t *testing.T) {
	rand.Seed(99)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		c    token.Credentials
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "first", c: token.Credentials{
			User: "pippo", Hashpass: hashstring.Md5Sum("pippo")},
			args: args{ctx: context.TODO()},
			want: "75ed1842-49e9-bc19-675e-4d1f766213da"},
		{name: "second", c: token.Credentials{
			User: "pippo", Hashpass: hashstring.Md5Sum("pipp")},
			args: args{ctx: context.TODO()},
			want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Token(tt.args.ctx); got != tt.want {
				t.Errorf("Credentials.token() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"context"
	"math/rand"
	"testing"
	"testtoken"

	"github.com/axamon/hashstring"
)

func Test_credentials_autenticato(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		c    testtoken.Credentials
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "first", c: testtoken.Credentials{User: "pippo", Hashpass: hashstring.Md5Sum("pippo")}, args: args{ctx: context.TODO()}, want: true},
		{name: "second", c: testtoken.Credentials{User: "pippo", Hashpass: hashstring.Md5Sum("pipp")}, args: args{ctx: context.TODO()}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.autenticato(tt.args.ctx); got != tt.want {
				t.Errorf("credentials.autenticato() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_credentials_token(t *testing.T) {
	rand.Seed(99)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		c    testtoken.Credentials
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "first", c: testtoken.Credentials{User: "pippo", Hashpass: hashstring.Md5Sum("pippo")}, args: args{ctx: context.TODO()}, want: "75ed1842-49e9-bc19-675e-4d1f766213da"},
		{name: "second", c: testtoken.Credentials{User: "pippo", Hashpass: hashstring.Md5Sum("pipp")}, args: args{ctx: context.TODO()}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.token(tt.args.ctx); got != tt.want {
				t.Errorf("credentials.token() = %v, want %v", got, tt.want)
			}
		})
	}
}

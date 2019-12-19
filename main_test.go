// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func Test_main(t *testing.T) {
	user = "pippo"
	pass = "pippo"
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{name: "first"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

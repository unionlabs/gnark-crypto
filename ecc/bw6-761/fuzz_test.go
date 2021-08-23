//go:build gofuzz
// +build gofuzz

// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bw6761

import (
	"encoding/hex"
	"io"
	"math/rand"
	"runtime/debug"
	"testing"
	"time"
)

func TestFuzz(t *testing.T) {
	const maxBytes = 1 << 10
	const testCount = 7
	var bytes [maxBytes]byte
	var i int
	seed := time.Now().UnixNano()
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
			t.Error(string(debug.Stack()))
			t.Fatal("test panicked", i, hex.EncodeToString(bytes[:i]), "seed", seed)
		}
	}()
	r := rand.New(rand.NewSource(seed))

	for i = 1; i < maxBytes; i++ {
		for j := 0; j < testCount; j++ {
			if _, err := io.ReadFull(r, bytes[:i]); err != nil {
				t.Fatal("couldn't read random bytes", err)
			}

			Fuzz(bytes[:i])
		}
	}

}

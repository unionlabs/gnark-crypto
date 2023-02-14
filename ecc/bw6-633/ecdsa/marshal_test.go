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

package ecdsa

import (
	"crypto/rand"
	"crypto/subtle"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

const (
	nbFuzzShort = 10
	nbFuzz      = 100
)

func TestSerialization(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[BW6-633] ECDSA serialization: SetBytes(Bytes()) should stay the same", prop.ForAll(
		func() bool {
			privKey, _ := GenerateKey(rand.Reader)

			var end PrivateKey
			buf := privKey.Bytes()
			n, err := end.SetBytes(buf[:])
			if err != nil {
				return false
			}
			if n != sizePrivateKey {
				return false
			}

			return end.PublicKey.Equal(&privKey.PublicKey) && subtle.ConstantTimeCompare(end.scalar[:], privKey.scalar[:]) == 1

		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

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

// Code generated by gurvy DO NOT EDIT

package fptower

// q (modulus)
var qE2 = [4]uint64{
	4332616871279656263,
	10917124144477883021,
	13281191951274694749,
	3486998266802970665,
}

// q'[0], see montgommery multiplication algorithm
var (
	qE2Inv0 uint64 = 9786893198990664585
	_              = qE2Inv0 // used in asm
)

//go:noescape
func addE2(res, x, y *E2)

//go:noescape
func subE2(res, x, y *E2)

//go:noescape
func doubleE2(res, x *E2)

//go:noescape
func negE2(res, x *E2)

//go:noescape
func mulNonResE2(res, x *E2)

//go:noescape
func squareAdxE2(res, x *E2)

//go:noescape
func mulAdxE2(res, x, y *E2)

// MulByNonResidue multiplies a E2 by (9,1)
func (z *E2) MulByNonResidue(x *E2) *E2 {
	mulNonResE2(z, x)
	return z
}

// Mul sets z to the E2-product of x,y, returns z
func (z *E2) Mul(x, y *E2) *E2 {
	mulAdxE2(z, x, y)
	return z
}

// Square sets z to the E2-product of x,x, returns z
func (z *E2) Square(x *E2) *E2 {
	squareAdxE2(z, x)
	return z
}

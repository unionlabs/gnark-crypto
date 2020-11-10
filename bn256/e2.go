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

package bn256

import (
	"math/big"

	"github.com/consensys/gurvy/bn256/fp"
)

// e2 is a degree two finite field extension of fp.Element
type e2 struct {
	A0, A1 fp.Element
}

// Equal returns true if z equals x, fasle otherwise
func (z *e2) Equal(x *e2) bool {
	return z.A0.Equal(&x.A0) && z.A1.Equal(&x.A1)
}

// Cmp compares (lexicographic order) z and x and returns:
//
//   -1 if z <  x
//    0 if z == x
//   +1 if z >  x
//
func (z *e2) Cmp(x *e2) int {
	if a1 := z.A1.Cmp(&x.A1); a1 != 0 {
		return a1
	}
	return z.A0.Cmp(&x.A0)
}

// LexicographicallyLargest returns true if this element is strictly lexicographically
// larger than its negation, false otherwise
func (z *e2) LexicographicallyLargest() bool {
	// adapted from github.com/zkcrypto/bls12_381
	if z.A1.IsZero() {
		return z.A0.LexicographicallyLargest()
	}
	return z.A1.LexicographicallyLargest()
}

// SetString sets a e2 element from strings
func (z *e2) SetString(s1, s2 string) *e2 {
	z.A0.SetString(s1)
	z.A1.SetString(s2)
	return z
}

// SetZero sets an e2 elmt to zero
func (z *e2) SetZero() *e2 {
	z.A0.SetZero()
	z.A1.SetZero()
	return z
}

// Set sets an e2 from x
func (z *e2) Set(x *e2) *e2 {
	z.A0 = x.A0
	z.A1 = x.A1
	return z
}

// SetOne sets z to 1 in Montgomery form and returns z
func (z *e2) SetOne() *e2 {
	z.A0.SetOne()
	z.A1.SetZero()
	return z
}

// SetRandom sets a0 and a1 to random values
func (z *e2) SetRandom() *e2 {
	z.A0.SetRandom()
	z.A1.SetRandom()
	return z
}

// IsZero returns true if the two elements are equal, fasle otherwise
func (z *e2) IsZero() bool {
	return z.A0.IsZero() && z.A1.IsZero()
}

// Add adds two elements of e2
func (z *e2) Add(x, y *e2) *e2 {
	addE2(z, x, y)
	return z
}

// Sub two elements of e2
func (z *e2) Sub(x, y *e2) *e2 {
	subE2(z, x, y)
	return z
}

// Double doubles an e2 element
func (z *e2) Double(x *e2) *e2 {
	doubleE2(z, x)
	return z
}

// Neg negates an e2 element
func (z *e2) Neg(x *e2) *e2 {
	negE2(z, x)
	return z
}

// String implements Stringer interface for fancy printing
func (z *e2) String() string {
	return (z.A0.String() + "+" + z.A1.String() + "*u")
}

// ToMont converts to mont form
func (z *e2) ToMont() *e2 {
	z.A0.ToMont()
	z.A1.ToMont()
	return z
}

// FromMont converts from mont form
func (z *e2) FromMont() *e2 {
	z.A0.FromMont()
	z.A1.FromMont()
	return z
}

// MulByElement multiplies an element in e2 by an element in fp
func (z *e2) MulByElement(x *e2, y *fp.Element) *e2 {
	var yCopy fp.Element
	yCopy.Set(y)
	z.A0.Mul(&x.A0, &yCopy)
	z.A1.Mul(&x.A1, &yCopy)
	return z
}

// Conjugate conjugates an element in e2
func (z *e2) Conjugate(x *e2) *e2 {
	z.A0 = x.A0
	z.A1.Neg(&x.A1)
	return z
}

// Legendre returns the Legendre symbol of z
func (z *e2) Legendre() int {
	var n fp.Element
	z.norm(&n)
	return n.Legendre()
}

// Exp sets z=x**e and returns it
func (z *e2) Exp(x e2, exponent *big.Int) *e2 {
	z.SetOne()
	b := exponent.Bytes()
	for i := 0; i < len(b); i++ {
		w := b[i]
		for j := 0; j < 8; j++ {
			z.Square(z)
			if (w & (0b10000000 >> j)) != 0 {
				z.Mul(z, &x)
			}
		}
	}

	return z
}

func init() {
	q := fp.Modulus()
	tmp := big.NewInt(3)
	sqrtExp1.Set(q).Sub(&sqrtExp1, tmp).Rsh(&sqrtExp1, 2)

	tmp.SetUint64(1)
	sqrtExp2.Set(q).Sub(&sqrtExp2, tmp).Rsh(&sqrtExp2, 1)
}

var sqrtExp1, sqrtExp2 big.Int

// Sqrt sets z to the square root of and returns z
// The function does not test wether the square root
// exists or not, it's up to the caller to call
// Legendre beforehand.
// cf https://eprint.iacr.org/2012/685.pdf (algo 9)
func (z *e2) Sqrt(x *e2) *e2 {

	var a1, alpha, b, x0, minusone e2

	minusone.SetOne().Neg(&minusone)

	a1.Exp(*x, &sqrtExp1)
	alpha.Square(&a1).
		Mul(&alpha, x)
	x0.Mul(x, &a1)
	if alpha.Equal(&minusone) {
		var c fp.Element
		c.Set(&x0.A0)
		z.A0.Neg(&x0.A1)
		z.A1.Set(&c)
		return z
	}
	a1.SetOne()
	b.Add(&a1, &alpha)

	b.Exp(b, &sqrtExp2).Mul(&x0, &b)
	z.Set(&b)
	return z
}

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

package pedersen

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Key for proof and verification
type Key struct {
	g             bn254.G2Affine // TODO @tabaie: does this really have to be randomized?
	gRootSigmaNeg bn254.G2Affine //gRootSigmaNeg = g^{-1/σ}
	basis         []bn254.G1Affine
	basisExpSigma []bn254.G1Affine
}

func randomOnG2() (bn254.G2Affine, error) { // TODO: Add to G2.go?
	gBytes := make([]byte, fr.Bytes)
	if _, err := rand.Read(gBytes); err != nil {
		return bn254.G2Affine{}, err
	}
	return bn254.HashToG2(gBytes, []byte("random on g2"))
}

func Setup(basis []bn254.G1Affine) (Key, error) {
	var (
		k   Key
		err error
	)

	if k.g, err = randomOnG2(); err != nil {
		return k, err
	}

	var modMinusOne big.Int
	modMinusOne.Sub(fr.Modulus(), big.NewInt(1))
	var sigma *big.Int
	if sigma, err = rand.Int(rand.Reader, &modMinusOne); err != nil {
		return k, err
	}
	sigma.Add(sigma, big.NewInt(1))

	var sigmaInvNeg big.Int
	sigmaInvNeg.ModInverse(sigma, fr.Modulus())
	sigmaInvNeg.Sub(fr.Modulus(), &sigmaInvNeg)
	k.gRootSigmaNeg.ScalarMultiplication(&k.g, &sigmaInvNeg)

	k.basisExpSigma = make([]bn254.G1Affine, len(basis))
	for i := range basis {
		k.basisExpSigma[i].ScalarMultiplication(&basis[i], sigma)
	}

	k.basis = basis
	return k, err
}

func (k *Key) Commit(values []fr.Element) (commitment bn254.G1Affine, knowledgeProof bn254.G1Affine, err error) {

	if len(values) != len(k.basis) {
		err = fmt.Errorf("unexpected number of values")
		return
	}

	// TODO @gbotrel this will spawn more than one task, see
	// https://github.com/ConsenSys/gnark-crypto/issues/269
	config := ecc.MultiExpConfig{
		NbTasks: 1, // TODO Experiment
	}

	if _, err = commitment.MultiExp(k.basis, values, config); err != nil {
		return
	}

	_, err = knowledgeProof.MultiExp(k.basisExpSigma, values, config)

	return
}

// VerifyKnowledgeProof checks if the proof of knowledge is valid
func (k *Key) VerifyKnowledgeProof(commitment bn254.G1Affine, knowledgeProof bn254.G1Affine) error {

	if !commitment.IsInSubGroup() || !knowledgeProof.IsInSubGroup() {
		return fmt.Errorf("subgroup check failed")
	}

	product, err := bn254.Pair([]bn254.G1Affine{commitment, knowledgeProof}, []bn254.G2Affine{k.g, k.gRootSigmaNeg})
	if err != nil {
		return err
	}
	if product.IsOne() {
		return nil
	}
	return fmt.Errorf("proof rejected")
}


func (k *Key) WriteTo(e *bn254.Encoder) (n int64, err error) {
	return k.writeTo(e)
}

func (k *Key) writeTo(e *bn254.Encoder) (int64, error) {
	if err := e.Encode(&k.g); err != nil {
		return e.BytesWritten(), err
	}
	if err := e.Encode(&k.gRootSigmaNeg); err != nil {
		return e.BytesWritten(), err
	}
	fmt.Printf("Basis         : %d\n", len(k.basis))
	if err := e.Encode(k.basis); err != nil {
		return e.BytesWritten(), err
	}
	fmt.Printf("BasisExpSignam: %d\n", len(k.basisExpSigma))
	if err := e.Encode(k.basisExpSigma); err != nil {
		return e.BytesWritten(), err
	}
	return e.BytesWritten(), nil
}

func (k *Key) ReadFrom(d *bn254.Decoder) (int64, error) {
	return k.readFrom(d)
}

func (k *Key) readFrom(d *bn254.Decoder) (int64, error) {
	if err := d.Decode(&k.g); err != nil {
		return d.BytesRead(), err
	}
	if err := d.Decode(&k.gRootSigmaNeg); err != nil {
		return d.BytesRead(), err
	}
	if err := d.Decode(&k.basis); err != nil {
		return d.BytesRead(), err
	}
	if err := d.Decode(&k.basisExpSigma); err != nil {
		return d.BytesRead(), err
	}
	return d.BytesRead(), nil
}

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Proof int64

// NewProof generates a proof based on the last one.
func NewProof(last Proof) (proof Proof) {
	for !proof.IsValid(last) {
		proof += 1
	}
	return proof
}

// IsValid checks that sha256 of last.proof contains 4 leading zeroes.
func (proof Proof) IsValid(last Proof) bool {
	sha := sha256.Sum256([]byte(fmt.Sprintf("%d.%d", last, proof)))
	return bytes.Compare(sha[:3], []byte{0, 0, 0}) == 0
}

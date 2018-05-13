package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Index int64
type Hash string

type Block struct {
	Index        Index
	Timestamp    int64
	Transactions []Transaction
	Proof        Proof
	Previous     Hash
}

// NewBlock creates a block in the block chain, given the proof.
func NewBlock(proof Proof) Block {
	var hash Hash
	var index Index
	if len(chain) > 0 {
		index = chain[len(chain)-1].Index + 1
		hash = chain[len(chain)-1].Hash()
	}

	block := Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		Proof:        proof,
		Previous:     hash,
	}
	chain = append(chain, block)

	// Reset current transactions.
	transactions = []Transaction{}

	return block
}

// Hash creates a SHA-256 hash of the Block.
func (b Block) Hash() Hash {
	sha := sha256.Sum256([]byte(fmt.Sprintf("%v", b)))
	return Hash(string(sha[:]))
}

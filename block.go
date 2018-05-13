package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Index int64
type Hash string

type Message struct {
	Sender, Receiver, Room string
	Text                   string
}

type Block struct {
	Index     Index
	Timestamp int64
	Message   Message
	Previous  Hash
}

func NewBlock(message Message) Block {
	var hash Hash
	var index Index
	if len(chain) > 0 {
		index = chain[len(chain)-1].Index + 1
		hash = chain[len(chain)-1].Hash()
	}

	block := Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		Message:   message,
		Previous:  hash,
	}
	chain = append(chain, block)

	return block
}

// Hash creates a SHA-256 hash of the Block.
func (b Block) Hash() Hash {
	sha := sha256.Sum256([]byte(fmt.Sprintf("%v", b)))
	return Hash(string(sha[:]))
}

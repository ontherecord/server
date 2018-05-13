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

func NewBlock(message Message) *Block {
	return &Block{
		Timestamp: time.Now().Unix(),
		Message:   message,
	}
}

// Hash creates a SHA-256 hash of the Block.
func (b Block) Hash() Hash {
	sha := sha256.Sum256([]byte(fmt.Sprintf("%v", b)))
	return Hash(string(sha[:]))
}

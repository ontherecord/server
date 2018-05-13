package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

type Chain struct {
	Id uuid.UUID

	// TODO: Maybe store as value copies to assert ownership.
	blocks []*Block
}

func NewChain() (c Chain) {
	c.Id = uuid.New()

	glog.Infof("Chain ID: %s", c.Id.String())

	// Create the genesis block and a node ID for this node.
	block := NewBlock(
		Message{
			Sender:   c.Id.String(),
			Receiver: c.Id.String(),
			Text:     "[genesis]",
		},
	)

	c.blocks = []*Block{block}
	return
}

func (c Chain) IsValid() bool {
	for i, b := range c.blocks {
		if b.Index != Index(i) {
			return false
		}
		if i == 0 {
			if b.Previous != "" {
				return false
			}
		} else if b.Previous != c.blocks[i-1].Hash() {
			return false
		}
	}
	return true
}

func (c Chain) Last() Block {
	return *c.blocks[len(c.blocks)-1]
}

func (c Chain) Add(block *Block) error {
	if len(c.blocks) == 0 {
		return fmt.Errorf("chain has not been initialized")
	}

	block.Index = c.Last().Index + 1
	block.Previous = c.Last().Hash()
	c.blocks = append(c.blocks, block)
	return nil
}

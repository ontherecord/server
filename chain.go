package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

type Chain struct {
	Id     uuid.UUID
	Blocks []Block
}

func NewChain() (c *Chain) {
	c = &Chain{Id: uuid.New()}

	glog.Infof("Chain ID: %s", c.Id.String())

	// Create the genesis block and a node ID for this node.
	block := NewBlock(
		Message{
			From: c.Id.String(),
			To:   c.Id.String(),
			Text: "[genesis]",
		},
	)

	c.Blocks = []Block{block}
	return
}

func (c Chain) IsValid() bool {
	for i, b := range c.Blocks {
		if b.Index != Index(i) {
			return false
		}
		if i == 0 {
			if b.Previous != "" {
				return false
			}
		} else if b.Previous != c.Blocks[i-1].Hash() {
			return false
		}
	}
	return true
}

func (c Chain) Last() Block {
	return c.Blocks[len(c.Blocks)-1]
}

func (c *Chain) Add(block Block) (Block, error) {
	if len(c.Blocks) == 0 {
		return block, fmt.Errorf("chain has not been initialized")
	}

	block.Index = c.Last().Index + 1
	block.Previous = c.Last().Hash()
	c.Blocks = append(c.Blocks, block)
	return block, nil
}

package main

import (
	"testing"
)

func TestNewBlock(t *testing.T) {
	cases := []struct {
		message Message
		want    Block
	}{}

	for _, tt := range cases {
		if got := NewBlock(tt.message); got != tt.want {
			t.Errorf("want %+v, got %+v", tt.want, got)
		}
	}
}

func TestHash(t *testing.T) {
	cases := []struct {
		block Block
		want  Hash
	}{}

	for _, tt := range cases {
		if got := tt.block.Hash(); got != tt.want {
			t.Errorf("want %q, got %q", tt.want, got)
		}
	}
}

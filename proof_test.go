package main

import "testing"

func TestNewProof(t *testing.T) {
	cases := []struct {
		last Proof
		want Proof
	}{
		{
			last: 42,
			want: 3316890,
		},
	}

	for _, tt := range cases {
		if got := NewProof(tt.last); got != tt.want {
			t.Errorf("want %d, got %d", tt.want, got)
		}
	}
}

func TestIsValid(t *testing.T) {
	cases := []struct {
		proof, last Proof
		want        bool
	}{
		{
			proof: 3316890,
			last:  42,
			want:  true,
		},
		{
			proof: 42,
			last:  42,
			want:  false,
		},
	}

	for _, tt := range cases {
		if got := tt.proof.IsValid(tt.last); got != tt.want {
			t.Errorf("want %t, got %t", tt.want, got)
		}
	}
}

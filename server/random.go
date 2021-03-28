package main

import (
	"github.com/zephyrtronium/crazy"

	"golang.org/x/exp/rand"
)

type MT64Source struct {
	mt64 *crazy.MT64
}

func (s MT64Source) Uint64() uint64 {
	return s.mt64.Uint64()
}

func (s MT64Source) Seed(seed uint64) {}

func NewCryptoSeededMT64() rand.Source {
	// NOTE: 312 is MT64's state's size
	return &MT64Source{crazy.CryptoSeeded(crazy.NewMT64(), 312).(*crazy.MT64)}
}

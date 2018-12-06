// Copyright (c) 2018 Australian Rivers Institute.

// Package rand implements an concurrency-safe wrapper around the language supplied rand package for
// pseudo-random number generators.  It does so by wrapping calls to an underlying math/rand Rand struct
// with a mutex lock.
package rand

import (
	"math/rand"
	"sync"
	"time"
)

// Container defines an interface embedding a Model
type Container interface {
	RandomNumberGenerator() *ConcurrencySafeRand
	SetRandomNumberGenerator(generator *ConcurrencySafeRand)
}

type ContainedRandomNumberGenerator struct {
	randomNumberGenerator *ConcurrencySafeRand
}

func (g *ContainedRandomNumberGenerator) RandomNumberGenerator() *ConcurrencySafeRand {
	return g.randomNumberGenerator
}

func (g *ContainedRandomNumberGenerator) SetRandomNumberGenerator(generator *ConcurrencySafeRand) {
	g.randomNumberGenerator = generator
}

// ConcurrencySafeRand is a concurrency-safe source of random numbers
type ConcurrencySafeRand struct {
	sync.Mutex
	unsafeRand *rand.Rand
}

// New returns a new ConcurrencySafeRand that uses random values from src to generate other random values.
func New(src rand.Source) *ConcurrencySafeRand {
	unsafeRand := rand.New(src)
	return &ConcurrencySafeRand{unsafeRand: unsafeRand}
}

// New returns a new ConcurrencySafeRand that uses random values seeded from a source of the system-time to generate
// other random values.
func NewTimeSeeded() *ConcurrencySafeRand {
	return New(rand.NewSource(time.Now().UnixNano()))
}

// Uint64 returns a pseudo-random 64-bit value as a uint64 from the default Source.
func (csr *ConcurrencySafeRand) Uint64() uint64 {
	csr.Lock()
	defer csr.Unlock()
	return csr.unsafeRand.Uint64()
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (csr *ConcurrencySafeRand) Intn(n int) int {
	csr.Lock()
	defer csr.Unlock()
	return csr.unsafeRand.Intn(n)
}

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func (csr *ConcurrencySafeRand) Int63n(n int64) int64 {
	csr.Lock()
	defer csr.Unlock()
	return csr.unsafeRand.Int63n(n)
}

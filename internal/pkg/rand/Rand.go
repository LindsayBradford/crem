// Copyright (c) 2018 Australian Rivers Institute.

// Package rand wraps the language supplied rand package for pseudo-random number generators.
// The wrapping supplies convenience functions for creation of typical seeding, and extra
// random functions not offered by the base library.

package rand

import (
	"math"
	"math/rand"
	"time"
)

// Container defines an interface embedding a Model
type Container interface {
	RandomNumberGenerator() *Rand
	SetRandomNumberGenerator(generator *Rand)
}

// RandContainer offers a struct implementing the Container interface.
type RandContainer struct {
	rand Rand
}

func (g *RandContainer) RandomNumberGenerator() *Rand {
	return &g.rand
}

func (g *RandContainer) SetRandomNumberGenerator(generator *Rand) {
	g.rand = *generator
}

// Rand is a source of project-specific random numbers
type Rand struct {
	officialRand rand.Rand
}

// New returns a new Rand that uses random values from src to generate other random values.
func New(src rand.Source) *Rand {
	unsafeRand := rand.New(src)
	return &Rand{officialRand: *unsafeRand}
}

// New returns a new Rand that uses random values seeded from a source of the system-time to generate
// other random values.
func NewTimeSeeded() *Rand {
	return New(rand.NewSource(time.Now().UnixNano()))
}

// Uint64 returns a pseudo-random 64-bit value as a uint64 from the default Source.
func (r *Rand) Uint64() uint64 {
	return r.officialRand.Uint64()
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (r *Rand) Intn(n int) int {
	return r.officialRand.Intn(n)
}

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func (r *Rand) Int63n(n int64) int64 {
	return r.officialRand.Int63n(n)
}

// Float64Unitary returns, as a float64, a non-negative pseudo-random number in [0,1].
func (r *Rand) Float64Unitary() float64 {
	// See: http://mumble.net/~campbell/2014/04/28/uniform-random-float (because it's tricky to get right)
	distributionRange := int64(math.Pow(2, 53))
	return float64(r.Int63n(distributionRange)) / float64(distributionRange-1)
}

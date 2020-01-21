package optimus

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

var (
	// MAX_INT represents the upper bound.
	// It should be 2^N.
	// It is set by default to the upper bound of an int32
	// to interface with https://github.com/jenssegers/optimus
	MAX_INT = uint64(math.MaxInt32) // 2,147,483,647
)

const (
	// MILLER_RABIN is used to configure the ProbablyPrime function
	// which is used to verify prime numbers.
	//
	// See: https://golang.org/pkg/math/big/#Int.ProbablyPrime
	MILLER_RABIN = 20
)

// Optimus is used to encode and decode integers using Knuth's Hashing Algorithm.
type Optimus struct {
	prime      uint64
	modInverse uint64
	random     uint64
}

// New returns an Optimus struct that can be used to encode and decode integers.
// A common use case is for obfuscating internal ids of database primary keys.
// It is imperative that you keep a record of prime, modInverse and random so that
// you can decode an encoded integer correctly. random must be an integer less than MAX_INT.
//
// WARNING: The function panics if prime is not a valid prime. It does a probability-based
// prime test using the MILLER-RABIN algorithm.
//
// CAUTION: DO NOT DIVULGE prime, modInverse and random!
func New(prime uint64, modInverse uint64, random uint64) Optimus {
	if prime > math.MaxInt64 {
		return Optimus{prime, modInverse, random}
	}

	p := big.NewInt(int64(prime))
	if p.ProbablyPrime(MILLER_RABIN) {
		return Optimus{prime, modInverse, random}
	}

	accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MILLER_RABIN))
	panic(fmt.Errorf("prime is not a valid prime. [Accuracy: %f]", accuracy))
}

// NewCalculated returns an Optimus struct that can be used to encode and decode integers.
// random must be an integer less than MAX_INT.
// It automatically calculates prime's mod inverse and then calls New.
func NewCalculated(prime uint64, random uint64) Optimus {
	return New(prime, ModInverse(prime), random)
}

// GenerateRandom generates a cryptographically secure random number.
// As currently implemented, it will return a number in the int64 range.
func GenerateRandom() uint64 {
	b49 := *big.NewInt(math.MaxInt64)
	n, _ := rand.Int(rand.Reader, &b49)
	in := n.Uint64() + 1

	return in
}

// Encode is used to encode n using Knuth's hashing algorithm.
func (o Optimus) Encode(n uint64) uint64 {
	return ((n * o.prime) & MAX_INT) ^ o.random
}

// Decode is used to decode n back to the original. It will only decode correctly if the Optimus struct
// is consistent with what was used to encode n.
func (o Optimus) Decode(n uint64) uint64 {
	return ((n ^ o.random) * o.modInverse) & MAX_INT
}

// Prime returns the associated prime.
//
// CAUTION: DO NOT DIVULGE THIS NUMBER!
func (o Optimus) Prime() uint64 {
	return o.prime
}

// ModInverse returns the associated mod inverse.
//
// CAUTION: DO NOT DIVULGE THIS NUMBER!
func (o Optimus) ModInverse() uint64 {
	return o.modInverse
}

// Random returns the associated random integer.
//
// CAUTION: DO NOT DIVULGE THIS NUMBER!
func (o Optimus) Random() uint64 {
	return o.random
}

// ModInverse returns the modular inverse of a given prime number.
// The modular inverse is defined such that
// (PRIME * MODULAR_INVERSE) & (MAX_INT_VALUE) = 1.
//
// See: http://en.wikipedia.org/wiki/Modular_multiplicative_inverse
//
// NOTE: prime is assumed to be a valid prime. If prime is outside the bounds of
// an int64, then the function panics as it can not calculate the mod inverse.
func ModInverse(prime uint64) uint64 {
	if prime > math.MaxInt64 {
		panic("prime exceeds max int64")
	}

	p := big.NewInt(int64(prime))

	var i big.Int
	max := big.NewInt(int64(MAX_INT + 1))

	return i.ModInverse(p, max).Uint64()
}

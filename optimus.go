package optimus

import (
	"fmt"
	"math"
	"math/big"

	"github.com/pjebs/jsonerror"
)

const (
	// MaxInt32 represents the maximum 32-bit integer allowed within the hashing
	// algorithm.
	MaxInt32 = 1<<31 - 1 // math.MaxInt32

	// MillerRabinRounds is the number of Miller-Rabin tests used to detect
	// primes.  See https://golang.org/pkg/math/big/#Int.ProbablyPrime
	MillerRabinRounds = 20
)

// Optimus represents the hashing implementation.
type Optimus struct {
	prime      uint64
	modInverse uint64
	random     uint64
}

// New returns an Optimus struct which can be used to encode and decode
// integers. Usually used for obfuscating internal ids such as database
// table rows.
func New(prime uint64, modInverse uint64, random uint64) Optimus {
	return Optimus{prime, modInverse, random}
}

// AssertPrime tests if a given number is prime.  If the test fails, a panic
// occurs.
func AssertPrime(n uint64) {
	p := big.NewInt(int64(n))
	if !p.ProbablyPrime(MillerRabinRounds) {
		accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MillerRabinRounds))
		panic(jsonerror.New(2, "Number is not prime", fmt.Sprintf("%d optimus.Miller-Rabin tests done. Accuracy: %f", MillerRabinRounds, accuracy)))
	}
}

// Encode encodes n using Knuth's Hashing Algorithm.
// Ensure that you store the prime, modInverse and random number
// associated with the Optimus struct so that it can be decoded
// correctly.
func (o Optimus) Encode(n uint64) uint64 {
	return ((n * o.prime) & MaxInt32) ^ o.random
}

// Decode decodes a number that had been hashed already using Knuth's Hashing Algorithm.
// It will only decode the number correctly if the prime, modInverse and random
// number associated with the Optimus struct is consistent with when the number
// was originally hashed.
func (o Optimus) Decode(n uint64) uint64 {
	return ((n ^ o.random) * o.modInverse) & MaxInt32
}

// Prime returns the Associated Prime Number. DO NOT DEVULGE THIS NUMBER!
func (o Optimus) Prime() uint64 {
	return o.prime
}

// ModInverse returns the Associated ModInverse Number. DO NOT DEVULGE THIS NUMBER!
func (o Optimus) ModInverse() uint64 {
	return o.modInverse
}

// Random returns the Associated Random Number. DO NOT DEVULGE THIS NUMBER!
func (o Optimus) Random() uint64 {
	return o.random
}

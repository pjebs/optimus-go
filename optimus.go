package optimus

import (
	"fmt"
	"math"
	"math/big"

	"github.com/pjebs/jsonerror"
)

const (
	MaxInt32          = 1<<31 - 1 // math.MaxInt32
	MillerRabinRounds = 20        //https://golang.org/pkg/math/big/#Int.ProbablyPrime
)

// Optimus represents the hashing implementation.
type Optimus struct {
	prime      uint64
	modInverse uint64
	random     uint64
}

// Returns an Optimus struct which can be used to encode and decode
// integers. Usually used for obfuscating internal ids such as database
// table rows. Panics if prime is not valid.
func New(prime uint64, modInverse uint64, random uint64) Optimus {
	p := big.NewInt(int64(prime))
	if !p.ProbablyPrime(MillerRabinRounds) {
		accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MillerRabinRounds))
		panic(jsonerror.New(2, "Number is not prime", fmt.Sprintf("%d Miller-Rabin tests done. Accuracy: %f", MillerRabinRounds, accuracy)))
	}

	return Optimus{prime, modInverse, random}
}

// Encodes n using Knuth's Hashing Algorithm.
// Ensure that you store the prime, modInverse and random number
// associated with the Optimus struct so that it can be decoded
// correctly.
func (this Optimus) Encode(n uint64) uint64 {
	return ((n * this.prime) & MaxInt32) ^ this.random
}

// Decodes a number that had been hashed already using Knuth's Hashing Algorithm.
// It will only decode the number correctly if the prime, modInverse and random
// number associated with the Optimus struct is consistent with when the number
// was originally hashed.
func (this Optimus) Decode(n uint64) uint64 {
	return ((n ^ this.random) * this.modInverse) & MaxInt32
}

// Returns the Associated Prime Number. DO NOT DEVULGE THIS NUMBER!
func (this Optimus) Prime() uint64 {
	return this.prime
}

// Returns the Associated ModInverse Number. DO NOT DEVULGE THIS NUMBER!
func (this Optimus) ModInverse() uint64 {
	return this.modInverse
}

// Returns the Associated Random Number. DO NOT DEVULGE THIS NUMBER!
func (this Optimus) Random() uint64 {
	return this.random
}

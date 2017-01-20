package optimus

import (
	"fmt"
	"math"
	"math/big"

	"github.com/pjebs/jsonerror"
)

var (
	MAX_INT = uint64(math.MaxInt32) //2147483647
)

const (
	MILLER_RABIN = 20 //https://golang.org/pkg/math/big/#Int.ProbablyPrime
)

type Optimus struct {
	prime      uint64
	modInverse uint64
	random     uint64
}

// Returns an Optimus struct which can be used to encode and decode
// integers. Usually used for obfuscating internal ids such as database
// table rows. Panics if prime is not valid.
// Warning: `ProbablyPrime` test is only done if `prime` is within bounds of `int64`
func New(prime uint64, modInverse uint64, random uint64) Optimus {
	if prime != uint64(int64(prime)) {
		return Optimus{prime, modInverse, random}
	}

	p := big.NewInt(int64(prime))
	if p.ProbablyPrime(MILLER_RABIN) {
		return Optimus{prime, modInverse, random}
	} else {
		accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MILLER_RABIN))
		panic(jsonerror.New(2, "Number is not prime", fmt.Sprintf("n=%d. %d Miller-Rabin tests done. Accuracy: %f", prime, MILLER_RABIN, accuracy)))
	}

}

// Returns an Optimus struct which can be used to encode and decode
// integers. Usually used for obfuscating internal ids such as database
// table rows. This method calculates the modInverse computationally.
// Panics if prime is not valid.
// Warning: `ProbablyPrime` test is only done if `prime` is within bounds of `int64`
func NewCalculated(prime uint64, random uint64) Optimus {
	prime64 := int64(prime)
	if prime != uint64(prime64) {
		return Optimus{prime, ModInverse(prime64), random}
	}

	p := big.NewInt(prime64)
	if p.ProbablyPrime(MILLER_RABIN) {
		return Optimus{prime, ModInverse(prime64), random}
	} else {
		accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MILLER_RABIN))
		panic(jsonerror.New(2, "Number is not prime", fmt.Sprintf("n=%d. %d Miller-Rabin tests done. Accuracy: %f", prime, MILLER_RABIN, accuracy)))
	}
}

// Encodes n using Knuth's Hashing Algorithm.
// Ensure that you store the prime, modInverse and random number
// associated with the Optimus struct so that it can be decoded
// correctly.
func (this Optimus) Encode(n uint64) uint64 {
	return ((n * this.prime) & MAX_INT) ^ this.random
}

// Decodes a number that had been hashed already using Knuth's Hashing Algorithm.
// It will only decode the number correctly if the prime, modInverse and random
// number associated with the Optimus struct is consistent with when the number
// was originally hashed.
func (this Optimus) Decode(n uint64) uint64 {
	return ((n ^ this.random) * this.modInverse) & MAX_INT
}

// Returns the Associated Prime Number. DO NOT DIVULGE THIS NUMBER!
func (this Optimus) Prime() uint64 {
	return this.prime
}

// Returns the Associated ModInverse Number. DO NOT DIVULGE THIS NUMBER!
func (this Optimus) ModInverse() uint64 {
	return this.modInverse
}

// Returns the Associated Random Number. DO NOT DIVULGE THIS NUMBER!
func (this Optimus) Random() uint64 {
	return this.random
}

// Calculates the Modular Inverse of a given Prime number such that
// (PRIME * MODULAR_INVERSE) & (MAX_INT_VALUE) = 1
// Panics if `prime` is not a valid prime number.
// See: http://en.wikipedia.org/wiki/Modular_multiplicative_inverse
func ModInverse(prime int64) uint64 {

	p := big.NewInt(prime)
	if !p.ProbablyPrime(MILLER_RABIN) {
		accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MILLER_RABIN))
		panic(jsonerror.New(2, "Number is not prime", fmt.Sprintf("n=%d. %d Miller-Rabin tests done. Accuracy: %f", prime, MILLER_RABIN, accuracy)))
	}

	var i big.Int
	max := big.NewInt(int64(MAX_INT + 1))

	return i.ModInverse(p, max).Uint64()
}

package optimus

import (
	"crypto/rand"
	"fmt"
	"github.com/caleblloyd/primesieve"
	"github.com/pjebs/jsonerror"
	"math"
	"math/big"
)

const (
	MAX_INT      = 2147483647
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
func New(prime uint64, modInverse uint64, random uint64) Optimus {

	p := big.NewInt(int64(prime))
	if p.ProbablyPrime(MILLER_RABIN) {
		return Optimus{prime, modInverse, random}
	} else {
		accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MILLER_RABIN))
		panic(jsonerror.New(2, "Number is not prime", fmt.Sprintf("%d Miller-Rabin tests done. Accuracy: %f", MILLER_RABIN, accuracy)))
	}

}

// Returns an Optimus struct which can be used to encode and decode
// integers. Usually used for obfuscating internal ids such as database
// table rows. This method calculates the modInverse computationally.
// Panics if prime is not valid.
func NewCalculated(prime uint64, random uint64) Optimus {
	p := big.NewInt(int64(prime))
	if p.ProbablyPrime(MILLER_RABIN) {
		return Optimus{prime, ModInverse(prime), random}
	} else {
		accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MILLER_RABIN))
		panic(jsonerror.New(2, "Number is not prime", fmt.Sprintf("%d Miller-Rabin tests done. Accuracy: %f", MILLER_RABIN, accuracy)))
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
// Panics if n is not a valid prime number.
// See: http://en.wikipedia.org/wiki/Modular_multiplicative_inverse
func ModInverse(n uint64) uint64 {
	p := big.NewInt(int64(n))

	if !p.ProbablyPrime(MILLER_RABIN) {
		accuracy := 1.0 - 1.0/math.Pow(float64(4), float64(MILLER_RABIN))
		panic(jsonerror.New(2, "Number is not prime", fmt.Sprintf("n=%d. %d Miller-Rabin tests done. Accuracy: %f", n, MILLER_RABIN, accuracy)))
	}

	var i big.Int

	prime := big.NewInt(int64(n))
	max := big.NewInt(int64(MAX_INT + 1))

	return i.ModInverse(prime, max).Uint64()
}

// Generates a valid Optimus struct by calculating a random prime number
// This function takes a few seconds and is resource intensive
func GenerateSeed() *Optimus {
	//Generated prime must be less than MAX_INT
	b := big.NewInt(MAX_INT - 1)
	n, _ := rand.Int(rand.Reader, b)
	//Calculates the largest prime less than n, this can take a few seconds
	selectedPrime := primesieve.PrimeMax(uint64(n.Int64()))

	//Calculate Mod Inverse for selectedPrime
	modInverse := ModInverse(selectedPrime)

	//Generate Random Integer less than MAX_INT
	upper := *big.NewInt(MAX_INT - 2)
	rand, _ := rand.Int(rand.Reader, &upper)
	randomNumber := rand.Uint64() + 1

	return &Optimus{selectedPrime, modInverse, randomNumber}
}

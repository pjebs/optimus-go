package optimus

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/pjebs/jsonerror"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"strings"
)

const (
	MAX_INT = 2147483647
)

type Optimus struct {
	prime      uint64
	modInverse uint64
	random     uint64
}

// Returns an Optimus struct which can be used to encode and decode
// integers. Usually used for obfuscating internal ids such as database
// table rows.
func New(prime uint64, modInverse uint64, random uint64) Optimus {
	return Optimus{prime, modInverse, random}
}

// Returns an Optimus struct which can be used to encode and decode
// integers. Usually used for obfuscating internal ids such as database
// table rows. This method calculates the modInverse computationally.
func NewCalculated(prime uint64, random uint64) Optimus {
	return Optimus{prime, ModInverse(prime), random}
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

// Calculates the Modular Inverse of a given Prime number such that
// (PRIME * MODULAR_INVERSE) & (MAX_INT_VALUE) = 1
// If n is not a Prime number, the return value is indeterminate.
// See: http://en.wikipedia.org/wiki/Modular_multiplicative_inverse
func ModInverse(n uint64) uint64 {
	var i big.Int

	prime := big.NewInt(int64(n))
	max := big.NewInt(int64(MAX_INT + 1))

	return i.ModInverse(prime, max).Uint64()
}

// Generates a valid Optimus struct using a randomly selected prime
// number from this site: http://primes.utm.edu/lists/small/millions/
// The first 50 million prime numbers are distributed evenly in 50 files.
// This Function is Time, Memory and CPU intensive. Run it once to generate the
// required seeds.
// WARNING: Potentially Unsecure. Double check that the prime number returned
// is actually prime number using an independent source.
// The largest Prime has 9 digits. The smallest has 1 digit.
// The final return value is the zipfile that was used to obtain the prime number
func GenerateSeed() (*Optimus, error, uint8) {
	log.Printf("WARNING: Optimus generates a random number via this site: http://primes.utm.edu/lists/small/millions/. This is potentially unsecure!")

	baseURL := "http://primes.utm.edu/lists/small/millions/primes%d.zip"

	//Generate Random number between 1-50
	b_49 := *big.NewInt(49)
	n, _ := rand.Int(rand.Reader, &b_49)
	i_n := n.Uint64() + 1

	//Download zip file
	finalUrl := fmt.Sprintf(baseURL, i_n)

	resp, err := client().Get(finalUrl)
	defer resp.Body.Close()
	if err != nil {
		return nil, jsonerror.New(1, "Could not generate seed", err.Error()), uint8(i_n)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, jsonerror.New(1, "Could not generate seed", err.Error()), uint8(i_n)
	}

	r, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
	if err != nil {
		return nil, jsonerror.New(1, "Could not generate seed", err.Error()), uint8(i_n)
	}

	zippedFile := r.File[0]

	src, err := zippedFile.Open() //src contains ReaderCloser
	defer src.Close()
	if err != nil {
		return nil, jsonerror.New(1, "Could not generate seed", err.Error()), uint8(i_n)
	}

	//Create a Byte Slice
	buf := new(bytes.Buffer)
	noOfBytes, _ := buf.ReadFrom(src)
	b := buf.Bytes() //Byte Slice

	//Randomly pick a character position
	start := 67 // Each zip file has an introductory header which is not relevant until the 67th character
	end := noOfBytes

	b_end := *big.NewInt(int64(end) - int64(start))
	n, _ = rand.Int(rand.Reader, &b_end)
	randomPosition := n.Uint64() + uint64(start)

	min := randomPosition - 9
	max := randomPosition + 9

	if min < uint64(start) {
		min = uint64(start)
	}

	if max > uint64(end) {
		max = uint64(end)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(b[min:max]))) //Input
	scanner.Split(bufio.ScanWords)

	var selectedNumbers []uint64
	for scanner.Scan() {
		p, _ := strconv.ParseUint(scanner.Text(), 10, 64)
		selectedNumbers = append(selectedNumbers, p)
	}

	//Not perfect but good enough

	var selectedPrime uint64
	length := len(selectedNumbers)
	if length > 2 {
		//Pick middle number

		//Check if length is even number
		//Check if round is odd or even
		var odd bool
		if length&1 != 0 {
			odd = true //odd
		} else {
			odd = false //even
		}

		if odd {
			selectedPrime = selectedNumbers[length/2]
		} else {

			r := *big.NewInt(1)
			rn, _ := rand.Int(rand.Reader, &r)
			if rn.Uint64() == 0 {
				selectedPrime = selectedNumbers[length/2]
			} else {
				selectedPrime = selectedNumbers[length/2-1]
			}
		}
	} else {
		//Pick largest number
		largest := selectedNumbers[0]

		for _, value := range selectedNumbers {
			if value > largest {
				largest = value
			}
		}

		selectedPrime = largest
	}

	//Calculate Mod Inverse for selectedPrime
	modInverse := ModInverse(selectedPrime)

	//Generate Random Integer less than MAX_INT
	upper := *big.NewInt(MAX_INT - 2)
	rand, _ := rand.Int(rand.Reader, &upper)
	randomNumber := rand.Uint64() + 1

	return &Optimus{selectedPrime, modInverse, randomNumber}, nil, uint8(i_n)
}

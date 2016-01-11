// +build !appengine

package optimus

import (
	"crypto/rand"
	"math/big"
	"testing"
)

// Obtains a prime number from the internet, calculates the mod inverse of it and
// calculates a random number. It then checks if the process worked BUT does not
// test if the number obtained is actually Prime.
func TestGenerateSeed(t *testing.T) {
	for i := 0; i < 3; i++ { //How many times we want to run GenerateSeed()
		o := GenerateSeed()
		t.Logf("Generated Seed - Prime: %d ModInverse: %d Random: %d", o.Prime(), o.ModInverse(), o.Random())
	}
}

// Tests if the encoding process correctly decodes the id back to the original.
func TestEncoding(t *testing.T) {
	for i := 0; i < 15; i++ { //How many times we want to run GenerateSeed()
		o := GenerateSeed()

		c := 10
		h := 100 //How many random numbers to select in between 0-c and (MAX_INT-c) - MAX-INT

		var y []uint64 //Stores all the values we want to run encoding tests on

		for t := 0; t < c; t++ {
			y = append(y, uint64(t))
		}

		//Generate Random numbers
		for t := 0; t < h; t++ {
			upper := *big.NewInt(int64(MAX_INT - 2*c))
			rand, _ := rand.Int(rand.Reader, &upper)
			randomNumber := rand.Uint64() + uint64(c)

			y = append(y, randomNumber)
		}

		for t := MAX_INT; t >= MAX_INT-c; t-- {
			y = append(y, uint64(t))
		}

		t.Logf("Prime: %d ModInverse: %d Random: %d", o.Prime(), o.ModInverse(), o.Random())
		for _, value := range y {
			orig := value
			hashed := o.Encode(value)
			unhashed := o.Decode(hashed)

			if orig != unhashed {
				t.Errorf("%d: %d -> %d - FAILED", orig, hashed, unhashed)
			} else {
				t.Logf("%d: %d -> %d - PASSED", orig, hashed, unhashed)
				// log.Printf("%d: %d -> %d - PASSED", orig, hashed, unhashed)
			}
		}

	}
}

package generator

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/pjebs/optimus-go"
)

// Tests if the encoding process correctly decodes the id back to the original.
func TestEncoding(t *testing.T) {

	for i := 0; i < 5; i++ { // How many times we want to run GenerateSeed()
		o, _ := GenerateSeed()

		c := 10
		h := 100 // How many random numbers to select in between 0-c and (MAX_INT-c) - MAX-INT

		var y []uint64 // Stores all the values we want to run encoding tests on

		for t := 0; t < c; t++ {
			y = append(y, uint64(t))
		}

		// Generate Random numbers
		for t := 0; t < h; t++ {
			upper := *big.NewInt(int64(optimus.MAX_INT - 2*uint64(c)))
			rand, _ := rand.Int(rand.Reader, &upper)
			randomNumber := rand.Uint64() + uint64(c)

			y = append(y, randomNumber)
		}

		for t := optimus.MAX_INT; t >= optimus.MAX_INT-uint64(c); t-- {
			y = append(y, t)
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
			}
		}

	}
}

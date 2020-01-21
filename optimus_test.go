package optimus

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

// Tests if the encoding process correctly decodes the id back to the original.
func TestEncoding(t *testing.T) {

	os := []Optimus{
		New(309779747, 49560203, 57733611),
		New(684934207, 1505143743, 846034763),
		New(743534599, 1356791223, 1336232185),
		New(54661037, 1342843941, 576322863),
		New(198194831, 229517423, 459462336),
	}

	for i := 0; i < 5; i++ { // How many times we want to run GenerateSeed()
		o := os[i]

		fmt.Println("o", o)

		c := 10
		h := 100 // How many random numbers to select in between 0-c and (MAX_INT-c) - MAX-INT

		var y []uint64 // Stores all the values we want to run encoding tests on

		for t := 0; t < c; t++ {
			y = append(y, uint64(t))
		}

		//Generate Random numbers
		for t := 0; t < h; t++ {
			upper := *big.NewInt(int64(MAX_INT - 2*uint64(c)))
			rand, _ := rand.Int(rand.Reader, &upper)
			randomNumber := rand.Uint64() + uint64(c)

			y = append(y, randomNumber)
		}

		for t := MAX_INT; t >= MAX_INT-uint64(c); t-- {
			y = append(y, t)
		}

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

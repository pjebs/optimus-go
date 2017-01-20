// +build !appengine

package generator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	// "log"
	"strings"
	"testing"
	"unsafe"

	"crypto/rand"
	// "log"
	"math/big"

	"github.com/pjebs/optimus-go"
)

// Obtains a prime number from the internet, calculates the mod inverse of it and
// calculates a random number. It then checks if the process worked BUT does not
// test if the number obtained is actually Prime.
func TestGenerateSeed(t *testing.T) {

	for i := 0; i < 3; i++ { //How many times we want to run GenerateSeed()
		o, f, err := GenerateSeed(nil)
		if err != nil {
			t.Errorf("Try %d - Failed", i)
		}

		//Check if prime is contained in zipped text file
		baseURL := "http://primes.utm.edu/lists/small/millions/primes%d.zip"
		finalUrl := fmt.Sprintf(baseURL, f)

		resp, err := client(nil).Get(finalUrl)
		if err != nil {
			t.Errorf("Try %d - Failed", i)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Try %d - Failed", i)
			continue
		}

		r, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
		if err != nil {
			t.Errorf("Try %d - Failed", i)
			continue
		}

		zippedFile := r.File[0]

		src, err := zippedFile.Open() //src contains ReaderCloser
		if err != nil {
			t.Errorf("Try %d - Failed", i)
			continue
		}
		// defer src.Close()

		//Create a Byte Slice
		buf := new(bytes.Buffer)
		buf.ReadFrom(src)
		b := buf.Bytes()
		stringContents := *(*string)(unsafe.Pointer(&b))

		subString := fmt.Sprintf(" %d ", o.Prime())
		if !strings.Contains(stringContents, subString) {
			src.Close()
			t.Errorf("Try %d - Failed - Obtained Prime:%d is not a Prime", i, o.Prime())
			continue
		}

		src.Close()

		//Check if ModInverse is correct

		// if o.Prime() != ModInverse(o.ModInverse()) {
		// 	src.Close()
		// 	t.Errorf("Try %d - Failed - ModInverse(%d) of %d is not correct", i, o.ModInverse, o.Prime)
		// 	continue
		// }

	}
}

// Tests if the encoding process correctly decodes the id back to the original.
func TestEncoding(t *testing.T) {

	for i := 0; i < 5; i++ { //How many times we want to run GenerateSeed()
		o, _, _ := GenerateSeed(nil)

		c := 10
		h := 100 //How many random numbers to select in between 0-c and (MAX_INT-c) - MAX-INT

		var y []uint64 //Stores all the values we want to run encoding tests on

		for t := 0; t < c; t++ {
			y = append(y, uint64(t))
		}

		//Generate Random numbers
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
				// log.Printf("%d: %d -> %d - PASSED", orig, hashed, unhashed)
			}
		}

	}
}

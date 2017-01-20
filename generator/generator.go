package generator

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/pjebs/jsonerror"
	"github.com/pjebs/optimus-go"
)

// Generates a valid Optimus struct using a randomly selected prime
// number from this site: http://primes.utm.edu/lists/small/millions/
// The first 50 million prime numbers are distributed evenly in 50 files.
// Parameter req should be nil if not using Google App Engine.
// This Function is Time, Memory and CPU intensive. Run it once to generate the
// required seeds.
// WARNING: Potentially Insecure. Double check that the prime number returned
// is actually prime number using an independent source.
// The largest Prime has 9 digits. The smallest has 1 digit.
// The second return value is the website zip file identifier that was used to obtain the prime number
func GenerateSeed(req *http.Request) (*optimus.Optimus, uint8, error) {
	log.Printf("\x1b[31mWARNING: Optimus generates a random number via this site: http://primes.utm.edu/lists/small/millions/. This is potentially insecure!\x1b[39;49m")

	baseURL := "http://primes.utm.edu/lists/small/millions/primes%d.zip"

	//Generate Random number between 1-50
	b_49 := *big.NewInt(49)
	n, _ := rand.Int(rand.Reader, &b_49)
	i_n := n.Uint64() + 1

	//Download zip file
	finalUrl := fmt.Sprintf(baseURL, i_n)
	log.Printf("Using file: %s", finalUrl)

	resp, err := client(req).Get(finalUrl)
	if err != nil {
		return nil, uint8(i_n), jsonerror.New(1, "Could not generate seed", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, uint8(i_n), jsonerror.New(1, "Could not generate seed", err.Error())
	}

	r, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
	if err != nil {
		return nil, uint8(i_n), jsonerror.New(1, "Could not generate seed", err.Error())
	}

	zippedFile := r.File[0]

	src, err := zippedFile.Open() //src contains ReaderCloser
	if err != nil {
		return nil, uint8(i_n), jsonerror.New(1, "Could not generate seed", err.Error())
	}
	defer src.Close()

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
	selectedPrime64 := int64(selectedPrime)
	if selectedPrime != uint64(selectedPrime64) {
		return nil, uint8(i_n), jsonerror.New(1, "Could not generate seed", "Prime number found by generator is too large to calculate the ModInverse. This is a limitation in math/big package. Try the generator again.")
	}
	modInverse := optimus.ModInverse(selectedPrime64)

	//Generate Random Integer less than MAX_INT
	upper := *big.NewInt(int64(optimus.MAX_INT - 2))
	rand, _ := rand.Int(rand.Reader, &upper)
	randomNumber := rand.Uint64() + 1

	o := optimus.New(selectedPrime, modInverse, randomNumber)
	return &o, uint8(i_n), nil
}

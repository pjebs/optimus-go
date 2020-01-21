package generator

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"

	"github.com/pjebs/optimus-go"
)

// GenerateSeed will generate a valid optimus object which can be used for encoding and
// decoding values.
//
// NOTE: You must download and place the list of prime numbers into directory "data".
// See: github.com/pjebs/optimus-go-primes (493+MB).
//
// If this repo is too large, then use RandN(50) to select a file from http://primes.utm.edu/lists/small/millions/
// and then Rand(1000000) to select a random prime.
func GenerateSeed() (*optimus.Optimus, error) {
	n := RandN(50)
	inputSource := fmt.Sprintf("./data/p%d.txt", n)
	lineNum := int(RandN(1000000))

	f, err := os.OpenFile(inputSource, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lineStr, err := realLine(f, lineNum)
	if err != nil {
		return nil, err
	}

	_selectedPrime, err := strconv.Atoi(lineStr)
	if err != nil {
		return nil, err
	}
	selectedPrime := uint64(_selectedPrime)
	modInverse := optimus.ModInverse(selectedPrime)
	random := RandN(int64(optimus.MAX_INT - 1))

	o := optimus.New(selectedPrime, modInverse, uint64(random))

	return &o, nil
}

func realLine(f io.Reader, lineNum int) (string, error) {

	lastLine := 0

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), nil
		}
	}
	err := sc.Err()
	if err != nil {
		return "", err
	}

	return "", nil
}

// RandN returns a cryptographically secure random number
// in the range [1,N].
func RandN(N int64) uint64 {
	b_49 := *big.NewInt(N)
	n, _ := rand.Int(rand.Reader, &b_49)
	i_n := n.Uint64() + 1
	return i_n
}

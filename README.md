ID Obfuscation/Hashing Transformer for Go [![GoDoc](http://godoc.org/github.com/pjebs/optimus-go?status.svg)](http://godoc.org/github.com/pjebs/optimus-go)
===============

There are many times when you want to generate obfuscated ids. This package utilizes Knuth's Hashing Algorithm to transform your internal ids into another number to *hide* it from the general public.

An example may be your database table. You may have a primary key that points to a particular customer. For security reasons you don't want to expose that id to the outside world. That is exactly where this package becomes handy.

Optimus encodes your internal id to a number that is safe to expose. Finally it can decode that number back so you know which internal id it refers to.

Full Unit Tests are provided.


Installation
-------------


```shell
go get -u github.com/pjebs/optimus-go
```


Usage
------

### Step 1

* Find or Calculate a **PRIME** number from [somewhere](http://primes.utm.edu/lists/small/millions/). It must be smaller than `2147483647` (MAXID)
* Calculate the Mod Inverse of the Prime number such that `(PRIME * INVERSE) & MAXID == 1`
* Generate a Pure Random Integer less than `2147483647` (MAXID).

You can use the built-in `GenerateSeed()` function to generate all 3 required parameters if you want. This is not recommended though because it is reliant on the prime numbers listed on the website: [http://primes.utm.edu/lists/small/millions/](http://primes.utm.edu/lists/small/millions/). This adds a point of point of insecurity.

If you do use the `GenerateSeed()` function, make sure that you verify:
* That website has not been hijacked
* The 50 million prime numbers listed on the site have not been modified to all point to the same number (i.e. they are all different)
* You independently verify that the Prime Number generated is in fact a **PRIME** number!


### Step 2

```go

package hello

import (
	"fmt"
	"github.com/pjebs/optimus"
	"net/http"
)

o, _, _ := optimus.New(1580030173, 59260789, 1163945558) //Prime Number: 1580030173, Mod Inverse: 59260789, Pure Random Number: 1163945558

new_id := o.Encode(15) //internal id of 15 being transformed

orig_id := o.Decode(1103647397) //Returns 15 back


```

Please note that in order for Optimus to transform the id back to the original, all 3 numbers of the constructor must be consistent. You will need to store it somewhere after generation and usage.

Methods
--------

```go
type Optimus struct {
	prime      uint64
	modInverse uint64
	random     uint64
}

```


```go
func New(prime uint64, modInverse uint64, random uint64) Optimus
```

Returns an Optimus struct which can be used to encode and decode integers. Usually used for obfuscating internal ids such as database table rows.


```go
func NewCalculated(prime uint64, random uint64) Optimus
```

Returns an Optimus struct which can be used to encode and decode integers. Usually used for obfuscating internal ids such as database table rows. This method calculates the modInverse computationally.

```go
func (this Optimus) Encode(n uint64) uint64 
```

Encodes n using Knuth's Hashing Algorithm.
Ensure that you store the prime, modInverse and random number associated with the Optimus struct so that it can be decoded correctly.

```go
func (this Optimus) Decode(n uint64) uint64
```

Decodes a number that had been hashed already using Knuth's Hashing Algorithm.
It will only decode the number correctly if the prime, modInverse and random number associated with the Optimus struct is consistent with when the number was originally hashed.

```go
func (this Optimus) Prime() uint64
```

Returns the Associated Prime Number. **DO NOT DEVULGE THIS NUMBER!**

```go
func (this Optimus) ModInverse() uint64
```

Returns the Associated ModInverse Number. **DO NOT DEVULGE THIS NUMBER!**

```go
func (this Optimus) Random() uint64
```

Returns the Associated Random Number. **DO NOT DEVULGE THIS NUMBER!**

```go
func ModInverse(n uint64) uint64
```

Calculates the Modular Inverse of a given Prime number such that `(PRIME * MODULAR_INVERSE) & (MAX_INT_VALUE) = 1`
If n is not a Prime number, the return value is indeterminate.
See: [http://en.wikipedia.org/wiki/Modular_multiplicative_inverse](http://en.wikipedia.org/wiki/Modular_multiplicative_inverse)

```go
func GenerateSeed() (*Optimus, error, uint8)
```

Generates a valid Optimus struct using a randomly selected prime number from this site: [http://primes.utm.edu/lists/small/millions/](http://primes.utm.edu/lists/small/millions/)
The first 50 million prime numbers are distributed evenly in 50 files.
This Function is Time, Memory and CPU intensive. Run it once to generate the required seeds.
**WARNING:** Potentially Insecure. Double check that the prime number returned is actually prime number using an independent source.
The largest Prime has 9 digits. The smallest has 1 digit.
The final return value is the website zip file identifier that was used to obtain the prime number

Alternatives
------------

There is the [hashids](http://hashids.org/) package which is very popular. Out of the box, it produces obfuscated ids that can contain any number of characters.

However:
* Knuth's algorithm is 127 times faster in benchmarks
* Hashids produce strings that contain characters other than just numbers.
	- If you were to modify the code (since the default *minimum* alphabet size is 16 characters) to allow only characters (0-9), it removes the first and last numbers to use as separators.
	- If the character '0' by coincidence comes out at the front of the obfuscated id, then you can't convert it to an integer when you store it. An integer will remove the leading zero but you need it to decode the number back to the original id (since hashid deals with strings and not numbers).

Inspiration
------------

This package is based on the PHP library by [jenssegers](https://github.com/jenssegers/optimus).

Final Notes
------------

If you found this package useful, please **Star** it on github. Feel free to fork or provide pull requests. Any bug reports will be warmly received.


[PJ Engineering and Business Solutions Pty. Ltd.](http://www.pjebs.com.au)
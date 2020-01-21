ID Obfuscation/Hashing Transformer for Go [![GoDoc](http://godoc.org/github.com/pjebs/optimus-go?status.svg)](http://godoc.org/github.com/pjebs/optimus-go)
===============

There are many times when you want to generate obfuscated ids. This package utilizes Knuth's Hashing Algorithm to transform your internal ids into another number to *hide* it from the general public.

An example may be your database table. You may have a primary key that points to a particular customer. For security reasons you don't want to expose that id to the outside world. That is exactly where this package becomes handy.

Optimus encodes your internal id to a number that is safe to expose. Finally it can decode that number back so you know which internal id it refers to.


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

You can use the built-in `generator.GenerateSeed()` function to generate all 3 required parameters if you want.


### Step 2

```go

package hello

import (
	"fmt"
	"github.com/pjebs/optimus-go"
)

o := optimus.New(1580030173, 59260789, 1163945558) // Prime Number: 1580030173, Mod Inverse: 59260789, Pure Random Number: 1163945558

new_id := o.Encode(15) // internal id of 15 being transformed to 1103647397

orig_id := o.Decode(1103647397) // Returns 15 back


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

New returns an Optimus struct that can be used to encode and decode integers.
A common use case is for obfuscating internal ids of database primary keys.
It is imperative that you keep a record of `prime`, `modInverse` and `random` so that
you can decode an encoded integer correctly. `random` must be an integer less than `MAX_INT`.

WARNING: The function panics if prime is not a valid prime. It does a probability-based
prime test using the MILLER-RABIN algorithm.

**CAUTION: DO NOT DIVULGE prime, modInverse and random!**


```go
func NewCalculated(prime uint64, random uint64) Optimus
```

NewCalculated returns an Optimus struct that can be used to encode and decode integers.
`random` must be an integer less than `MAX_INT`.
It automatically calculates prime's mod inverse and then calls New.


```go
func (this Optimus) Encode(n uint64) uint64 
```

Encode is used to encode n using Knuth's hashing algorithm.

```go
func (this Optimus) Decode(n uint64) uint64
```

Decode is used to decode n back to the original. It will only decode correctly if the Optimus struct is consistent with what was used to encode n.

```go
func (this Optimus) Prime() uint64
```

Prime returns the associated prime.

**CAUTION: DO NOT DIVULGE THIS NUMBER!**

```go
func (this Optimus) ModInverse() uint64
```

ModInverse returns the associated mod inverse.

**CAUTION: DO NOT DIVULGE THIS NUMBER!**

```go
func (this Optimus) Random() uint64
```

Random returns the associated random integer.

**CAUTION: DO NOT DIVULGE THIS NUMBER!**

```go
func ModInverse(n int64) uint64
```


ModInverse returns the modular inverse of a given prime number.
The modular inverse is defined such that `(PRIME * MODULAR_INVERSE) & (MAX_INT_VALUE) = 1`.

See: http://en.wikipedia.org/wiki/Modular_multiplicative_inverse

NOTE: prime is assumed to be a valid prime. If prime is outside the bounds of
an int64, then the function panics as it can not calculate the mod inverse.



```go
func generator.GenerateSeed() (*Optimus, error)
```

GenerateSeed will generate a valid optimus object which can be used for encoding and decoding values.

See http://godoc.org/github.com/pjebs/optimus-go/generator#GenerateSeed for details on how to use it.

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

If you found this package useful, please **Star** it on github. Feel free to fork and/or provide pull requests. Any bug reports will be warmly received.

**© 2014-20 PJ Engineering and Business Solutions Pty. Ltd.**

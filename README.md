pnc
=========

A library of pseudorandom number generators (PRNGs), ideal for playing with cryptography. These are not built to be secure and **will get you hacked if deployed in a real system**.

**pnc** contains standalone, clean implementations. You should be able to:

* **Generate** numbers as part of a larger system.
* **Attack** third-party generators to recover their internal state.
* **Understand** and modify the code to better understand how the generators work.

## Usage
Install [Golang](https://golang.org/doc/install) on your machine.

Run `go get github.com/46bit/pnc` to fetch the code.

Try the examples in `$GOPATH/github.com/46bit/pnc/examples`.

## Generators implemented
* **Mersenne Twister** [[1](https://en.wikipedia.org/wiki/Mersenne_twister), [2](http://www.quadibloc.com/crypto/co4814.htm)], an incredibly common, insecure PRNG.
* **Linear Congruential Generator** [[1](https://en.wikipedia.org/wiki/Linear_congruential_generator)], a common but defective PRNG.
* **Tausworthe Generator** [[1](http://www.cs.rice.edu/~johnmc/comp528/lecture-notes/Lecture21.pdf), [2](http://www.randombit.net/bitbashing/2008/07/01/linux_random32_failure_case.html)], a weak PRNG once used in slot machines.
* **Blum Blum Shub** [[1](https://en.wikipedia.org/wiki/Blum_Blum_Shub), [2](https://www.princeton.edu/~achaney/tmve/wiki100k/docs/Blum_Blum_Shub.html)], a CSPRNG secured by integer factorisation.
* **Blum Micali** [[1](https://en.wikipedia.org/wiki/Blum%E2%80%93Micali_algorithm)], a CSPRNG secured by the discrete-logarithm problem (DLP).
* **Dual EC DRBG** [[1](http://blog.cryptographyengineering.com/2013/09/the-many-flaws-of-dualecdrbg.html), [2](http://www.untruth.org/~josh/school/phd/seminar/spring-2013-dual-ec-drbg/DualECDRBG.pdf), [3](https://projectbullrun.org/dual-ec/)], a famously backdoored CSPRNG secured by elliptic curve DLP.

## About
Built by [Michael Mokrysz](https://46b.it) from December 2013. Licensed under MIT and formerly known as Pinocchio.

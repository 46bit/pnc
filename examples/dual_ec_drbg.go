package main

import (
  "fmt"
  "math/big"
  "github.com/46bit/pinocchio"
  "github.com/46bit/pinocchio/ec"
)

// NIST Curve P-256 (http://csrc.nist.gov/groups/ST/toolkit/documents/dss/NISTReCur.pdf)
// Q sourced from http://csrc.nist.gov/publications/nistpubs/800-90A/SP800-90A.pdf
const (
  curve_p256_p = "115792089210356248762697446949407573530086143415290314195533631308867097853951"
  curve_p256_n = "115792089210356248762697446949407573529996955224135760342422259061068512044369"
  curve_p256_a = "-3"
  curve_p256_b = "5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b"
  curve_p256_h = "1"
  curve_p256_gx = "6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296"
  curve_p256_gy = "4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5"
  dual_ec_drbg_curve_p256_qx = "c97445f45cdef9f0d3e05e1e585fc297235b82b5be8ff3efca67c59852018192"
  dual_ec_drbg_curve_p256_qy = "b28ef557ba31dfcbdd21ac46e2a91e3c304f44cb87058ada2cb815151e610046"
)

func main() {
  curve := ec.NewPrimeCurve(
    ec.NewBigInt(curve_p256_p, 10),
    ec.NewBigInt(curve_p256_a, 10),
    ec.NewBigInt(curve_p256_b, 16),
    ec.NewBigInt(curve_p256_gx, 16),
    ec.NewBigInt(curve_p256_gy, 16),
    ec.NewBigInt(curve_p256_n, 10),
    ec.NewBigInt(curve_p256_h, 10))

  // ----------------

  // Check generator point is on the curve.
  if !curve.Satisfied(curve.G) {
    fmt.Println("Generator point not on curve!")
  }

  // ----------------

  // Check a specified point is on the curve.
  p := ec.NewPoint(
    "84658269074130531148357510928537398853078465142270397446803060672548601752264",
    "32478479853643897696398775628369298920280560255617340607511813418971428614048",
    10)
  if !curve.Satisfied(p) {
    fmt.Println("Known-good point not on curve!")
  }

  // ----------------

  // Check add routine seems to keep points on the curve (a *weak* suggestion
  // the routine is correct).
  if !curve.Satisfied(curve.Add(p, curve.G)) {
    fmt.Println("Generator point plus known-good point not on curve!")
  }

  // ----------------

  // Check doubling routines seem to keep points on the curve (a *weak* suggestion
  // the routine is correct).
  p = curve.G
  for i := 1; i <= 100; i++ {
    p = curve.Double(p)
    if !curve.Satisfied(p) {
      fmt.Printf("Generator point doubled %d times is not on curve!\n", i)
    }
    if !p.Finite {
      fmt.Printf("Generator point doubled %d times was infinite!\n", i)
    }
  }

  // ----------------

  p2 := &ec.Point{big.NewInt(0), big.NewInt(0), false}
  p2 = curve.Add(curve.G, p2)
  p2 = curve.Double(p2)
  if curve.Satisfied(p2) {
    fmt.Println("Operations upon infinite point were on curve!")
  }
  if p2.Finite {
    fmt.Println("Operations upon infinite point were finite!")
  }

  // ----------------

  // Check scalar multiplication routine seems to keep points on the curve (a *weak*
  // suggestion the routine is correct).
  p = curve.G
  n := big.NewInt(0)
  n.SetString("e43fe3e4729612e27820e8c9ae246baf43f1d4281fcdc2ed", 16)
  p2 = curve.ScalarMultiply(n, p)
  if !curve.Satisfied(p2) {
    fmt.Printf("Generator Point Multiplied < N was not on curve!", n)
  }
  if !p2.Finite {
    fmt.Printf("Generator Point Multiplied < N was infinite!", n)
  }

  // ----------------

  // Curious case of buggy point addition.

  // This addition works with r1.x < r2.x despite point addition being commutative
  // in nature. WTF, suggests something broken despite many known good results.
  p1 := ec.NewPoint(
    "48439561293906451759052585252797914202762949526041747995844080717082404635286",
    "36134250956749795798585127919587881956611106672985015071877198253568414405109",
    10)

  p2 = ec.NewPoint(
    "56515219790691171413109057904011688695424810155802929973526481321309856242040",
    "3377031843712258259223711451491452598088675519751548567112458094635497583569",
    10)

  p3 := curve.Add(p1, p2)
  if !curve.Satisfied(p3) {
    fmt.Println("Non-commutative points where p1.x < p2.x was not on curve!")
  }
  if !p3.Finite {
    fmt.Println("Non-commutative point where p1.x < p2.x was infinite!")
  }

  p4 := curve.Add(p2, p1)
  if !curve.Satisfied(p3) {
    fmt.Println("Non-commutative point where p1.x > p2.x was not on curve!")
  }
  if !p3.Finite {
    fmt.Println("Non-commutative point where p1.x > p2.x was infinite!")
  }

  if !p3.Eq(p4) {
    fmt.Println("Non-commutative points were not equal when swapped!")
  }

  // ----------------

  // Generate 20 32-bit random integers using Dual_EC_DRBG on NIST Curve-256.
  // NB: There is every reason to suspect the NSA can recover internal PRNG state
  // when using the standard curves, and output demonstrates known biases using ~any
  // curve. Plus, even a better optimized version is particularly slow.
  s := ec.NewBigInt("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
  g := pinocchio.NewDualECDRBG(
    curve,
    ec.NewBigInt(dual_ec_drbg_curve_p256_qx, 16),
    ec.NewBigInt(dual_ec_drbg_curve_p256_qy, 16),
    s)
  for i := 0; i < 20000; i++ {
    fmt.Printf("%d: %d\n", i, g.Urand32())
  }

  // ----------------
}

package main

import "fmt"
import "math/big"
import "github.com/46bit/pinocchio"

// Curve P-256 (http://csrc.nist.gov/publications/drafts/800-90/draft_sp800_90a_rev1.pdf)
// Q sourced from http://csrc.nist.gov/publications/nistpubs/800-90A/SP800-90A.pdf
const (
  dual_ec_drbg_n =  "115792089210356248762697446949407573529996955224135760342422259061068512044369"
  dual_ec_drbg_a =  "-3"
  dual_ec_drbg_b =  "5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b"
  dual_ec_drbg_fp = "115792089210356248762697446949407573530086143415290314195533631308867097853951"
  dual_ec_drbg_px = "6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296"
  dual_ec_drbg_py = "4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5"
  dual_ec_drbg_qx = "c97445f45cdef9f0d3e05e1e585fc297235b82b5be8ff3efca67c59852018192"
  dual_ec_drbg_qy = "b28ef557ba31dfcbdd21ac46e2a91e3c304f44cb87058ada2cb815151e610046"
)

func main() {
  curve := pinocchio.NewECCurve(dual_ec_drbg_n, dual_ec_drbg_a, dual_ec_drbg_b, dual_ec_drbg_fp, dual_ec_drbg_px, dual_ec_drbg_py)

  // ----------------

  // Check curve point P is on the curve.
  if curve.Satisfied(curve.P) {
    fmt.Println("curve.P on curve :)")
  } else {
    fmt.Println("curve.P not on curve :(")
  }
  fmt.Println("---")

  // ----------------

  // Check a specified point is on the curve.
  o_p := pinocchio.NewECPoint(
    "84658269074130531148357510928537398853078465142270397446803060672548601752264",
    "32478479853643897696398775628369298920280560255617340607511813418971428614048",
    10)
  if curve.Satisfied(o_p) {
    fmt.Println("o_p on curve :)")
  } else {
    fmt.Println("o_p not on curve :(")
  }
  fmt.Println("---")

  // ----------------

  // Check doubling routines seem to keep points on the curve (a *weak* suggestion
  // the routine is correct).
  q := curve.P
  for i := 1; i <= 10; i++ {
    q = curve.Double(q)
    if curve.Satisfied(q) {
      fmt.Printf("curve.Q*2(%d) on curve :)\n", i)
    } else {
      fmt.Printf("curve.Q*2(%d) not on curve :(\n", i)
    }
    fmt.Println("---")
  }
  fmt.Println("---")

  // ----------------

  // Check add routine seems to keep points on the curve (a *weak* suggestion
  // the routine is correct).
  if curve.Satisfied(curve.Add(curve.P, curve.P)) {
    fmt.Println("curve.Q+curve.Q on curve :)")
  } else {
    fmt.Println("curve.Q+curve.Q not on curve :(")
  }
  fmt.Println("---")

  // ----------------

  // @TODO: Need to handle curve points at infinity. Example not implemented for now.
  r2 := &pinocchio.ECPoint{big.NewInt(0), big.NewInt(0), false}
  r2 = curve.Add(curve.P, r2)
  r2 = curve.Double(r2)
  if curve.Satisfied(r2) {
    fmt.Println("curve.P+curve.Q on curve :)")
  } else {
    fmt.Println("curve.P+curve.Q not on curve :(")
  }
  fmt.Println("---")

  // ----------------

  // Check scalar multiplication routine seems to keep points on the curve (a *weak*
  // suggestion the routine is correct).
  p := curve.P
  n := big.NewInt(0)
  n.SetString("1", 10)
  if curve.Satisfied(curve.ScalarMultiply(n, p)) {
    fmt.Printf("curve.Q*%d on curve :)\n", n)
  } else {
    fmt.Printf("curve.Q*%d not on curve :(\n", n)
  }
  fmt.Println("---")

  // ----------------

  // Curious case of buggy point addition.

  // This addition works with r1.x < r2.x despite point addition being commutative
  // in nature. WTF, suggests something broken despite many known good results.
  r2 = pinocchio.NewECPoint(
    "56515219790691171413109057904011688695424810155802929973526481321309856242040",
    "3377031843712258259223711451491452598088675519751548567112458094635497583569",
    10)

  r1 := pinocchio.NewECPoint(
    "48439561293906451759052585252797914202762949526041747995844080717082404635286",
    "36134250956749795798585127919587881956611106672985015071877198253568414405109",
    10)

  r3 := curve.Add(r1, r2)
  if curve.Satisfied(r3) {
    fmt.Println("curve.P+curve.Q on curve :)")
  } else {
    fmt.Println("curve.p+curve.Q not on curve :(")
  }
  fmt.Println("---")

  // ----------------

  // Generate 20 32-bit random integers using Dual_EC_DRBG on NIST Curve-256.
  // NB: There is every reason to suspect the NSA can recover internal PRNG state
  // when using the standard curves, and output demonstrates known biases using ~any
  // curve. Plus, even a better optimized version is particularly slow.
  s := big.NewInt(0)
  s.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
  g := pinocchio.NewDualECDRBG(curve, dual_ec_drbg_qx, dual_ec_drbg_qy, s)
  for i := 0; i < 20; i++ {
    fmt.Printf("%d: %d\n", i, g.Urand32())
  }
  fmt.Println("---")
}

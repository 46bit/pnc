package main

import (
  "fmt"
  "github.com/46bit/pnc"
  "github.com/46bit/pnc/ec"
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

  // Generate pseudorandom bytes using Dual_EC_DRBG on NIST Curve-256.
  // NB: Never, ever use this generator. It is ridiculously slow, demonstrates bias
  // and for the provided values of Q is backdoored by the NSA.

  // The seed s is the value of S *after* seeding the OpenSSL implementation.
  // Any integer on the order of 2^256 will suffice.
  // @TODO: Have compatible seeding routines with OpenSSL.
  s := ec.NewBigInt("14611F02F7F34E6121433EFB0D71ECAC38F28BE4274B3DD784D2C1D4BE78DF89", 16)
  g := pnc.NewDualECDRBG(
    curve,
    ec.NewBigInt(dual_ec_drbg_curve_p256_qx, 16),
    ec.NewBigInt(dual_ec_drbg_curve_p256_qy, 16),
    s)

  for i := 0; i < 55; i++ {
    fmt.Printf("%d: %x\n", i, g.Byte())
  }
}

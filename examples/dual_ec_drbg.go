package main

import (
  "fmt"
  "github.com/46bit/pnc"
  "github.com/46bit/pnc/ec"
)

// NIST Curve P-256 (http://csrc.nist.gov/groups/ST/toolkit/documents/dss/NISTReCur.pdf)
// Q sourced from http://csrc.nist.gov/publications/nistpubs/800-90A/SP800-90A.pdf
const (
  dual_ec_drbg_curve_p256_qx = "c97445f45cdef9f0d3e05e1e585fc297235b82b5be8ff3efca67c59852018192"
  dual_ec_drbg_curve_p256_qy = "b28ef557ba31dfcbdd21ac46e2a91e3c304f44cb87058ada2cb815151e610046"
)

func main() {
  // Generate pseudorandom bytes using Dual_EC_DRBG on NIST Curve-256.
  // NB: Never, ever use this generator. It is ridiculously slow, demonstrates bias
  // and for the provided values of Q is backdoored by the NSA.

  // The seed s is the value of S *after* seeding the OpenSSL implementation.
  // Any integer on the order of 2^256 will suffice.
  // @TODO: Have compatible seeding routines with OpenSSL.
  s := ec.NewBigInt("14611F02F7F34E6121433EFB0D71ECAC38F28BE4274B3DD784D2C1D4BE78DF89", 16)

  curve := ec.NewP256Curve()
  g := pnc.NewDualECDRBG(
    curve,
    ec.NewBigInt(dual_ec_drbg_curve_p256_qx, 16),
    ec.NewBigInt(dual_ec_drbg_curve_p256_qy, 16),
    s)

  for i := 0; i < 10; i++ {
    fmt.Printf("%x", g.Bytes(600))
  }
  fmt.Println()
}

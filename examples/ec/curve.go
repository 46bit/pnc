package main

import (
  "fmt"
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

  p := curve.G.Copy()
  fmt.Println("Before:")
  p.Print()
  r := curve.Add(p, p)
  fmt.Println("After:")
  r.Print()
}

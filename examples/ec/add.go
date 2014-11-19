package main

import (
  "fmt"
  "math/big"
  "github.com/46bit/pnc/ec"
)

// NIST Curve P-256 (http://csrc.nist.gov/groups/ST/toolkit/documents/dss/NISTReCur.pdf)
// Q sourced from http://csrc.nist.gov/publications/nistpubs/800-90A/SP800-90A.pdf
const (
  dual_ec_drbg_curve_p256_qx = "c97445f45cdef9f0d3e05e1e585fc297235b82b5be8ff3efca67c59852018192"
  dual_ec_drbg_curve_p256_qy = "b28ef557ba31dfcbdd21ac46e2a91e3c304f44cb87058ada2cb815151e610046"
)

func main() {
  curve := ec.NewP256Curve()

  p := curve.G.Copy()
  q := ec.NewPoint(dual_ec_drbg_curve_p256_qx, dual_ec_drbg_curve_p256_qy, 16)

  t := big.NewInt(0)
  t.SetString("05ABA71EB402603B7D24D9F921E49433A69AB3DB2D5A9910FF040FA906207587", 16)

  // ----------------

  fmt.Println("Before:")
  p.Print()

  for i := 0; i < 1931; i++ {
    p = curve.Add(p, q)
  }

  fmt.Println("After:")
  p.Print()
}

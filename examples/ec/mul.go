package main

import (
  "fmt"
  "math/big"
  "github.com/46bit/pnc/ec"
)

func main() {
  curve := ec.NewP256Curve()

  p := curve.G.Copy()

  t := big.NewInt(0)
  t.SetString("05ABA71EB402603B7D24D9F921E49433A69AB3DB2D5A9910FF040FA906207587", 16)

  // ----------------

  fmt.Println("Before:")
  p.Print()

  r := curve.ScalarMultiply(t, p)

  fmt.Println("After:")
  r.Print()

  fmt.Printf("Expected:\n- x = %X\n- y = %X\n- on curve: %t\n",
             ec.NewBigInt("7FDA41915769256A2D8F968BC9897849FC44C5CA64CF03E576EAF95E5FF9A799", 16),
             ec.NewBigInt("D7E013E76E4CEDCEB49F8C267164954F0D57C3FD077B0A81DF4DDA5AF4D5868D", 16),
             true)
}

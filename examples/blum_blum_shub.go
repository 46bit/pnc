package main

import (
  "fmt"
  "math/big"
  "github.com/46bit/pnc"
)

func main() {
  // @TODO: These p, q, s are far too small, with a series length of only 12.
  p := big.NewInt(11)
  q := big.NewInt(19)
  s := big.NewInt(3)

  b := pnc.NewBlumBlumShub()
  b.Seed(p, q, s)
  fmt.Printf("p = %d, q = %d, s = %d", p, q, s)

  for i := 0; i < 20; i++ {
    v := b.Urand32()
    fmt.Printf("%dth value is %d\n", i, v)
  }
}

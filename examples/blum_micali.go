package main

import (
  "fmt"
  "math/big"
  "github.com/46bit/pnc"
)

func main() {
  p := big.NewInt(1169939)
  g := big.NewInt(69937)
  s := big.NewInt(3)

  b := pnc.NewBlumMicali()
  b.Seed(p, g, s)

  for i := 0; i < 20; i++ {
    fmt.Printf("%d: %d\n", i, b.Urand32())
  }
}

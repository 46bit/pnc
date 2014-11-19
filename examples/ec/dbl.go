package main

import (
  "fmt"
  "github.com/46bit/pnc/ec"
)

func main() {
  curve := ec.NewP256Curve()

  p := curve.G.Copy()

  // ----------------

  fmt.Println("Before:")
  p.Print()

  for i := 0; i < 10000; i++ {
    p = curve.Double(p)
  }

  fmt.Println("After:")
  p.Print()
}

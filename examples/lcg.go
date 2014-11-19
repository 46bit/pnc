package main

import (
  "fmt"
  "github.com/46bit/pnc"
)

func main() {
  l := pnc.NewLCG()
  l.Seed(0)

  for i := 0; i < 20; i++ {
    fmt.Printf("%d: %d\n", i, l.Urand32())
  }
}

package main

import (
  "os"
  "fmt"
  "log"
  "bufio"
  "github.com/46bit/pnc"
)

func main() {
  // Seed a mersenne twister.
  g := pnc.NewMersenneTwister(3754129273)

  // Print and write 624 integers to a file.
  file, err := os.Create("/tmp/mersenne_outputs.txt")
  if err != nil {
    log.Fatalf("Couldn't create file: %s", err)
  }
  defer file.Close()
  w := bufio.NewWriter(file)

  for i := 0; i < 624; i++ {
    r := g.Urand32()
    fmt.Printf("#%d: %d\n", i, r)
    fmt.Fprintln(w, r)
  }

  w.Flush()
}

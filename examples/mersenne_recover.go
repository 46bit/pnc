package main

import (
  "os"
  "fmt"
  "log"
  "bufio"
  "strconv"
  "github.com/46bit/pnc"
)

func main() {
  g := pnc.NewMersenneTwister(3754129273)

  // Read sequence of 624 integers from a mersenne twister.
  // Try by running examples/mersenne.go then this.
  file, err := os.Open("/tmp/mersenne_outputs.txt")
  if err != nil {
    log.Fatalf("Couldn't read file: %s", err)
  }
  defer file.Close()
  s := bufio.NewScanner(file)

  var urand32s [624]uint32
  i := 0
  for s.Scan() {
    urand64, err := strconv.ParseUint(s.Text(), 16, 0)
    if err != nil {
      log.Fatalf("Error converting line %d to an integer: %s", i, err)
    }
    urand32s[i] = uint32(urand64) // inputs are 32-bit, this loses nothing
    i++
  }
  if i != 624 {
    log.Fatalf("Read %d integers. Need 624 from a mersenne twister.", i)
  }

  // Recover the state of the mersenne twister.
  m := pnc.NewMersenneTwister(0)
  m.SeedFromUrand32s(urand32s)

  // Print the next 624 integers of the mersenne twister.
  for i := 0; i < 624; i++ {
    fmt.Printf("#%d: %d\n", i, g.Urand32())
  }
}

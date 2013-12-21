package main

import (
  "fmt"
)

// with the restriction that 2nw − r − 1 is a Mersenne prime
const (
  w = iota // word size (in number of bits)
  n = iota // degree of recurrence
  m = iota // middle word, or the number of parallel sequences, 1 ≤ m ≤ n
  r = iota // separation point of one word, or the number of bits of the lower bitmask, 0 ≤ r ≤ w - 1
  a = iota // coefficients of the rational normal form twist matrix
  b = iota // TGFSR(R) tempering bitmasks
  c = iota // TGFSR(R) tempering bitmasks
  s = iota // TGFSR(R) tempering bit shifts
  t = iota // TGFSR(R) tempering bit shifts
  u = iota // additional Mersenne Twister tempering bit shifts
  l = iota // additional Mersenne Twister tempering bit shifts
)

type MersenneTwister struct {
  State [624]uint32
  index int
}

func NewMersenneTwister() MersenneTwister {
  m := MersenneTwister{}
  return m
}

func (m *MersenneTwister) Seed(seed uint32) {
  m.index = 0
  m.State[0] = seed
  for i := 1; i < 624; i++ {
    t := m.State[i - 1] ^ (m.State[i - 1] >> 30)
    m.State[i] = 1812433253 * (t + uint32(i))
  }
}

func (m *MersenneTwister) generate_numbers() {
  for i := 0; i < 624; i++ {
    y := (m.State[i] & 0x80000000) + (m.State[(i + 1) % 624] & 0x7fffffff)
    m.State[i] = m.State[(i + 397) % 624] ^ (y >> 1)
    if y % 2 != 0 {
      m.State[i] = m.State[i] ^ 2567483615
    }
  }
}

func (m *MersenneTwister) Urand32() uint32 {
  if m.index == 0 {
    m.generate_numbers()
  }

  y := m.State[m.index]
  y = y ^ (y >> 11)
  y = y ^ ((y << 7) & 2636928640)
  y = y ^ ((y << 15) & 4022730752)
  y = y ^ (y >> 18)

  m.index = (m.index + 1) % 624
  return y
}

func main() {
  debug := true

  m := NewMersenneTwister()
  m.Seed(0)

  for i := 0; r < 1000000; i++ {
    r := m.Urand32()
    fmt.Printf("%dth value is %d\n", i, r)
  }

  if debug { fmt.Println("EOF") }
}

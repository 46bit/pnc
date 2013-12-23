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

func (m *MersenneTwister) coded_xor(urand32 uint32, changebit_index uint32, xorbit_index uint32) uint32 {
  changebit_shift_index := 31 - changebit_index
  xorbit_shift_index := 31 - xorbit_index

  changebit := (urand32 >> changebit_shift_index) & 1
  xorbit := (urand32 >> xorbit_shift_index) & 1
  result := changebit ^ xorbit

  if changebit == 1 && result == 0 {
    urand32 -= (1 << changebit_shift_index)
  } else if changebit == 0 && result == 1 {
    urand32 += (1 << changebit_shift_index)
  }

  return urand32
}

func (m *MersenneTwister) coded_and_xor(urand32 uint32, uand32 uint32, changebit_index uint32, xorbit_index uint32) uint32 {
  xorbit_shift_index := 31 - xorbit_index
  changebit_shift_index := 31 - changebit_index

  changebit := (urand32 >> changebit_shift_index) & 1
  uandbit := (uand32 >> changebit_shift_index) & 1
  xorbit := (urand32 >> xorbit_shift_index) & 1
  result := xorbit & uandbit

  result = changebit ^ result

  if changebit == 1 && result == 0 {
    urand32 -= (1 << changebit_shift_index)
  } else if changebit == 0 && result == 1 {
    urand32 += (1 << changebit_shift_index)
  }

  return urand32
}

func (m *MersenneTwister) Urand32ToState(urand32 uint32) uint32 {
  for i := uint32(18); i <= 31; i++ {
    urand32 = m.coded_xor(urand32, i, i - 18)
  }

  for i := uint32(16); true; i-- {
    urand32 = m.coded_and_xor(urand32, 4022730752, i, i + 15)
    if i == 0 { break }
  }

  for i := uint32(24); true; i-- {
    urand32 = m.coded_and_xor(urand32, 2636928640, i, i + 7)
    if i == 0 { break }
  }

  for i := uint32(11); i <= 31; i++ {
    urand32 = m.coded_xor(urand32, i, i - 11)
  }

  return urand32
}

func main() {
  m := NewMersenneTwister()
  m.Seed(0)

  for i := 0; i < 623; i++ {
    v := m.Urand32()

    li := m.index - 1
    if li < 0 { li = 0 }
    s0 := m.State[li]
    s1 := m.Urand32ToState(v)

    if s0 != s1 {
      fmt.Printf("%d %d, %d != %d\n", i, v, s0, s1)
    }
  }
}

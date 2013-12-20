package main

import (
//  "os"
  "fmt"
//  "log"
//  "time"
//  "math/rand"
//  "bytes"
//  "bufio"
//  "io/ioutil"
//  "regexp"
//  "strings"
//  "strconv"
//  "encoding/json"
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

type Mt struct {
  state [624]int
  index int
}

func NewMt(seed int) Mt {
  m := Mt{}
  m.index = 0
  m.initialize_generator(seed)
  return m
}

func (m *Mt) initialize_generator(seed int) {
  m.state[0] = seed
  for i := 1; i < 624; i++ {
    t := m.state[i - 1] ^ (m.state[i - 1] >> 30)
    m.state[i] = 1812433253 * (t + i)
  }
}

func (m *Mt) generate_numbers() {
  for i := 0; i < 624; i++ {
    y := (m.state[i] & 0x80000000) + (m.state[(i + 1) % 624] & 0x7fffffff)
    m.state[i] = m.state[(i + 397) % 624] ^ (y >> 1)
    if y % 2 != 0 {
      m.state[i] = m.state[i] ^ 2567483615
    }
  }
}

func (m *Mt) extract_number() int {
  if m.index == 0 {
    m.generate_numbers()
  }

  y := m.state[m.index]
  y = y ^ (y >> 11)
  y = y ^ ((y << 7) & 2636928640)
  y = y ^ ((y << 15) & 4022730752)
  y = y ^ (y >> 18)

  m.index = (m.index + 1) % 624
  return y
}

func main() {
  debug := true

  m := NewMt(0)
  //for i := 0; i < 1000000; i++ {
  //r := 0
  for r, i := 0, 0; r != 1; i++ {
    r = m.extract_number()
    fmt.Printf("%dth value is %d\n", i, r)
  }

  if debug { fmt.Println("EOF") }
}

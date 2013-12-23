package main

import (
  "fmt"
  "math/big"
)

type BlumMicali struct {
  P *big.Int
  G *big.Int
  Term *big.Int
}

func NewBlumMicali() BlumMicali {
  b := BlumMicali{}
  return b
}

// p: large prime
// g: primitive root modulo p
func (b *BlumMicali) Seed(p *big.Int, g *big.Int, s *big.Int) {
  b.P = p
  b.G = g
  b.Term = s
}

func (b *BlumMicali) generate_next_term() {
  b.Term.Exp(b.G, b.Term, b.P)
}

func (b *BlumMicali) Urand32() uint32 {
  urand32 := uint32(0)
  for i := 0; i < 32; i++ {
    b.generate_next_term()

    numerator := big.NewInt(1)
    numerator.Sub(b.P, numerator)
    denominator := big.NewInt(2)

    n := big.NewInt(0)
    n.Div(numerator, denominator)

    urand32 = urand32 << 1
    if b.Term.Cmp(n) == -1 {
      urand32 += 1
    }
  }
  return urand32
}

func main() {
  p := big.NewInt(1169939)
  g := big.NewInt(69937)
  s := big.NewInt(3)

  b := NewBlumMicali()
  b.Seed(p, g, s)

  a := b.Urand32()
  //for i := 0; i < 50; i++ {
  for i := 0; true; i++ {
    v := b.Urand32()
    //if i % 10000 == 0 {
      fmt.Printf("%dth value No Match %d\n", i, v)
    //}
    if a == v {
      fmt.Printf("%dth value is %d\n", 1, v)
      //break
    }
  }
}

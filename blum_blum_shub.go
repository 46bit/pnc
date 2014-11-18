package pnc

import (
  "math/big"
)

type BlumBlumShub struct {
  P *big.Int
  Q *big.Int
  M *big.Int
  Term *big.Int
}

func NewBlumBlumShub() BlumBlumShub {
  b := BlumBlumShub{}
  return b
}

// p: large prime congruent to 3 mod 4, with small gcd(φ(p-1), φ(q-1))
// q: large prime congruent to 3 mod 4, with small gcd(φ(p-1), φ(q-1))
// s: integer coprime to M=PQ, not 0 or 1
func (b *BlumBlumShub) Seed(p *big.Int, q *big.Int, s *big.Int) {
  b.P = p
  b.Q = q
  b.Term = s

  b.M = big.NewInt(0)
  b.M.Mul(b.P, b.Q)
}

func (b *BlumBlumShub) generate_next_term() {
  b.Term.Exp(b.Term, big.NewInt(2), b.M)
}

func (b *BlumBlumShub) Bit() uint32 {
  b.generate_next_term()
  return uint32(b.Term.Bit(0))
}

// We generate uint32 from the LSB of 32 terms. Therefore the
// periodicity really needs to be a large multiple of that.
// @TODO: determine bits we can extract using http://www.win.tue.nl/~berry/papers/ima05bbs.pdf
// @TODO: refactor PRNGs to output specified n bits as opposed to uint32.
func (b *BlumBlumShub) Urand32() uint32 {
  urand32 := uint32(0)
  for i := 0; i < 32; i++ {
    urand32 = urand32 << 1 + b.Bit()
  }
  return urand32
}

/*
func main() {
  // @TODO: These p, q, s are far too small, with a series length of only 12.
  p := big.NewInt(11)
  q := big.NewInt(19)
  s := big.NewInt(3)

  b := NewBlumBlumShub()
  b.Seed(p, q, s)
  fmt.Printf("p = %d, q = %d, s = %d", p, q, s)

  for i := 0; i < 20; i++ {
    v := b.Urand32()
    fmt.Printf("%dth value is %d\n", i, v)
  }
}
*/

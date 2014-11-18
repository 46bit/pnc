package ec

import (
  "math/big"
)

// Implementation of elliptic curve on prime field operations.
// As with the rest of PNC, this is a rough toolkit to
// break things rather than something to ever deploy. This suffers
// from severe side channel attacks and I'll probably build a demo
// of that one day.

func NewBigInt(v string, base int) *big.Int {
  b := big.NewInt(0)
  b.SetString(v, base)
  return b
}

type Point struct {
  X *big.Int
  Y *big.Int
  Finite bool
}

func NewPoint(x string, y string, base int) *Point {
  return &Point{NewBigInt(x, base), NewBigInt(y, base), true}
}

func InfinitePoint() *Point {
  return &Point{big.NewInt(0), big.NewInt(0), false}
}

func (p *Point) Eq(p2 *Point) bool {
  return p.X.String() == p2.X.String() && p.Y.String() == p2.Y.String() && p.Finite == p2.Finite
}

func (p *Point) Copy() *Point {
  p2 := NewPoint(p.X.String(), p.Y.String(), 10)
  p2.Finite = p.Finite
  return p2
}

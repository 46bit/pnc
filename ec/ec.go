package ec

import (
  "math/big"
)

// Implementation of elliptic curve on prime field operations.
// As with the rest of Pinocchio, this is a rough toolkit to
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
  Infinite bool
}

func NewPoint(x string, y string, base int) *Point {
  p := &Point{NewBigInt(x, base), NewBigInt(y, base), false}
  // @TODO: This check will cause problems with initialising points as blank.
  //p.Infinite = p.X.Cmp(NewBigInt(0, 10)) == 0
  return p
}

func (o *Point) Copy() *Point {
  n := NewPoint(o.X.String(), o.Y.String(), 10)
  n.Infinite = o.Infinite
  return n
}

type PrimeCurve struct {
  P *big.Int // Prime modulus
  A *big.Int // y^2 = x^3 + Ax + b
  B *big.Int // y^2 = x^3 + ax + B
  G *Point // Generator
  N *big.Int // Order of Generator
  H *big.Int // Cofactor
}

func NewPrimeCurve(p, a, b, gx, gy, n, h *big.Int) *PrimeCurve {
  c := PrimeCurve{p, a, b, &Point{gx, gy, false}, n, h}
  return &c
}

func (c *PrimeCurve) Satisfied(p *Point) bool {
  // On curve if (y^2) - (x^3 + ax + b (mod p)) == 0

  // y^2
  y2 := big.NewInt(0)
  y2.Exp(p.Y, big.NewInt(2), nil)
  y2.Mod(y2, c.P)

  // x^3
  x3 := big.NewInt(0)
  x3.Exp(p.X, big.NewInt(3), nil)

  // ax
  ax := big.NewInt(0)
  ax.Mul(p.X, c.A)

  // x^3 + ax + b (mod p)
  rhs := big.NewInt(0)
  rhs.Add(x3, ax)
  rhs.Add(rhs, c.B)
  rhs.Mod(rhs, c.P)

  // (y^3 (mod p)) - (x^3 + ax + b (mod p))
  diff := big.NewInt(0)
  diff.Sub(y2, rhs)

  return diff.Cmp(big.NewInt(0)) == 0
}

func (c *PrimeCurve) Add(p1 *Point, p2 *Point) *Point {
  // Handle points adding to 0 (Infinite).

  if p1.X.String() == p2.X.String() && p1.Y.String() == p2.Y.String() {
    return c.Double(p1)
  }

  // Addition must be commutative. Yet this routine works IFF p1.X < p2.X.
  // I don't understand this, and that should scare you as much as it does me.
  // Nothing I've read says this constraint should be necessary.
  // TL;DR hey something's wrong!
  // @TODO: Work out if/where this function still fails and see if you can
  // identify the mistaken math. Number Theory class may help.
  if p1.X.Cmp(p2.X) > 0 {
    tp := p1.Copy()
    p1 = p2
    p2 = tp
  }

  st := big.NewInt(0)
  st.Sub(p2.Y, p1.Y)
  sb := big.NewInt(0)
  sb.Sub(p2.X, p1.X)

  // Multiply by multiplicative inverse, not integer division.
  mi := big.NewInt(0)
  mi.ModInverse(sb, c.P)
  st.Mul(st, mi)

  p3 := Point{big.NewInt(0), big.NewInt(0), false}

  p3.X.Exp(st, big.NewInt(2), nil)
  p3.X.Sub(p3.X, p1.X)
  p3.X.Sub(p3.X, p2.X)

  p3.Y.Sub(p1.X, p3.X)
  p3.Y.Mul(st, p3.Y)
  p3.Y.Sub(p3.Y, p1.Y)

  p3.X.Mod(p3.X, c.P)
  p3.Y.Mod(p3.Y, c.P)

  return &p3
}

func (c *PrimeCurve) Double(p1 *Point) *Point {
  s := big.NewInt(0)
  s.Exp(p1.X, big.NewInt(2), nil)
  s.Mul(s, big.NewInt(3))
  s.Add(s, c.A)
  k := big.NewInt(0)
  k.Mul(big.NewInt(2), p1.Y)

  // Multiply by multiplicative inverse, not integer division.
  mi := big.NewInt(0)
  mi.ModInverse(k, c.P)
  s.Mul(s, mi)

  p2 := &Point{big.NewInt(0), big.NewInt(0), false}

  p2.X.Exp(s, big.NewInt(2), nil)
  px2 := big.NewInt(0)
  px2.Mul(big.NewInt(2), p1.X)
  p2.X.Sub(p2.X, px2)

  p2.Y.Sub(p1.X, p2.X)
  p2.Y.Mul(p2.Y, s)
  p2.Y.Sub(p2.Y, p1.Y)

  p2.X.Mod(p2.X, c.P)
  p2.Y.Mod(p2.Y, c.P)

  return p2
}

func (c *PrimeCurve) ScalarMultiply(scalar *big.Int, p1 *Point) *Point {
  //scalar.Mod(scalar, curve.P)

  r := p1.Copy()

  for i := 0; i < scalar.BitLen(); i++ {
    if scalar.Bit(i) == 1 {
      r = c.Add(r, p1)
    }
    r = c.Double(r)
  }

  return r
}

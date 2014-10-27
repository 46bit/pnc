package pinocchio

import (
  "math/big"
)

// Implementation of elliptic curve on prime field operations.
// As with the rest of Pinocchio, this is a rough toolkit to
// break things rather than something to ever deploy. This suffers
// from severe side channel attacks and I'll probably build a demo
// of that one day.

type ECPoint struct {
  X *big.Int
  Y *big.Int
  Infinite bool
}

func NewECPoint(x string, y string, base int) *ECPoint {
  new_point := ECPoint{big.NewInt(0), big.NewInt(0), false}
  new_point.X.SetString(x, base)
  new_point.Y.SetString(y, base)
  return &new_point
}

func (old_point *ECPoint) CopyECPoint() *ECPoint {
  new_point := NewECPoint(old_point.X.String(), old_point.Y.String(), 10)
  new_point.Infinite = old_point.Infinite
  return new_point
}

type ECCurve struct {
  N *big.Int
  A *big.Int
  B *big.Int
  Fp *big.Int
  P *ECPoint
}

func NewECCurve(n, a, b, fp, px, py string) *ECCurve {
  curve := ECCurve{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil}

  curve.N.SetString(n, 10)
  curve.A.SetString(a, 10)
  curve.B.SetString(b, 16)
  curve.Fp.SetString(fp, 10)

  curve.P = NewECPoint(px, py, 16)

  return &curve
}

func (curve *ECCurve) Satisfied(p *ECPoint) bool {
  // On curve if (y^2) - (x^3 + ax + b (mod p)) == 0

  // y^2
  y2 := big.NewInt(0)
  y2.Exp(p.Y, big.NewInt(2), nil)
  y2.Mod(y2, curve.Fp)

  // x^3
  x3 := big.NewInt(0)
  x3.Exp(p.X, big.NewInt(3), nil)

  // ax
  ax := big.NewInt(0)
  ax.Mul(p.X, curve.A)

  // x^3 + ax + b (mod p)
  rhs := big.NewInt(0)
  rhs.Add(x3, ax)
  rhs.Add(rhs, curve.B)
  rhs.Mod(rhs, curve.Fp)

  // (y^3 (mod p)) - (x^3 + ax + b (mod p))
  diff := big.NewInt(0)
  diff.Sub(y2, rhs)

  return diff.Cmp(big.NewInt(0)) == 0
}

func (curve *ECCurve) Add(p1 *ECPoint, p2 *ECPoint) *ECPoint {
  if p1.X.String() == p2.X.String() && p1.Y.String() == p2.Y.String() {
    return curve.Double(p1)
  }

  // Addition must be commutative. Yet this routine works IFF p1.X < p2.X.
  // I don't understand this, and that should scare you as much as it does me.
  // Nothing I've read says this constraint should be necessary.
  // TL;DR hey something's wrong!
  // @TODO: Work out if/where this function still fails and see if you can
  // identify the mistaken math. Number Theory class may help.
  if p1.X.Cmp(p2.X) > 0 {
    tp := p1.CopyECPoint()
    p1 = p2
    p2 = tp
  }

  st := big.NewInt(0)
  st.Sub(p2.Y, p1.Y)
  sb := big.NewInt(0)
  sb.Sub(p2.X, p1.X)

  // Multiply by multiplicative inverse, not integer division.
  mi := big.NewInt(0)
  mi.ModInverse(sb, curve.Fp)
  st.Mul(st, mi)

  p3 := ECPoint{big.NewInt(0), big.NewInt(0), false}

  p3.X.Exp(st, big.NewInt(2), nil)
  p3.X.Sub(p3.X, p1.X)
  p3.X.Sub(p3.X, p2.X)

  p3.Y.Sub(p1.X, p3.X)
  p3.Y.Mul(st, p3.Y)
  p3.Y.Sub(p3.Y, p1.Y)

  p3.X.Mod(p3.X, curve.Fp)
  p3.Y.Mod(p3.Y, curve.Fp)

  return &p3
}

func (curve *ECCurve) Double(p1 *ECPoint) *ECPoint {
  p2 := ECPoint{big.NewInt(0), big.NewInt(0), false}
  s := big.NewInt(0)
  s.Exp(p1.X, big.NewInt(2), nil)
  s.Mul(s, big.NewInt(3))
  s.Add(s, curve.A)
  k := big.NewInt(0)
  k.Mul(big.NewInt(2), p1.Y)

  // Multiply by multiplicative inverse, not integer division.
  mi := big.NewInt(0)
  mi.ModInverse(k, curve.Fp)
  s.Mul(s, mi)

  p2.X.Exp(s, big.NewInt(2), nil)
  px2 := big.NewInt(0)
  px2.Mul(big.NewInt(2), p1.X)
  p2.X.Sub(p2.X, px2)

  p2.Y.Sub(p1.X, p2.X)
  p2.Y.Mul(p2.Y, s)
  p2.Y.Sub(p2.Y, p1.Y)

  p2.X.Mod(p2.X, curve.Fp)
  p2.Y.Mod(p2.Y, curve.Fp)

  return &p2
}

func (curve *ECCurve) ScalarMultiply(scalar *big.Int, p1 *ECPoint) *ECPoint {
  //scalar.Mod(scalar, curve.Fp)

  r := p1.CopyECPoint()

  for i := 0; i < scalar.BitLen(); i++ {
    if scalar.Bit(i) == 1 {
      r = curve.Add(r, p1)
    }
    r = curve.Double(r)
  }

  return r
}

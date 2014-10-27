package pinocchio

import (
  "fmt"
  "math/big"
)

type DualECDRBG struct {
  curve *ECCurve
  Q *ECPoint

  Z *big.Int
  S *big.Int
  StateIndex uint64
  StateBit uint32
}

func NewDualECDRBG(curve *ECCurve, qx, qy string, seed *big.Int) *DualECDRBG {
  g := DualECDRBG{curve, NewECPoint(qx, qy, 16), big.NewInt(0), big.NewInt(0), 0, 0}
  g.seed(seed)
  return &g
}

// Seed is twice security_strength bits long (at least 256)
func (g *DualECDRBG) seed(seed *big.Int) {
  g.S = seed
  g.generate_number()
}

func (g *DualECDRBG) generate_number() {
  g.StateIndex++
  g.StateBit = 0
  s := g.curve.ScalarMultiply(g.S, g.curve.P)
  if !g.curve.Satisfied(s) {
    fmt.Println("s not on curve")
  }
  z := g.curve.ScalarMultiply(s.X, g.Q)
  if !g.curve.Satisfied(z) {
    fmt.Println("z not on curve")
  }
  g.S = s.X
  g.Z = z.X
}

func (g *DualECDRBG) Bit() uint32 {
  bit := g.Z.Bit(int(g.StateBit))
  g.StateBit++
  if int(g.StateBit) >= g.Z.BitLen() - 16 {
    g.generate_number()
  }
  return uint32(bit)
}

func (g *DualECDRBG) Urand32() uint32 {
  v := uint32(0)
  for i := 0; i < 32; i++ {
    v = (v<<1) + g.Bit()
  }
  return v
}

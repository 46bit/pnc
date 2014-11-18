package pinocchio

import (
  "errors"
  "math/big"
  "github.com/46bit/pinocchio/ec"
)

type DualECDRBG struct {
  C *ec.PrimeCurve
  Q *ec.Point

  S *big.Int
  Sp *ec.Point
  Z *big.Int
  Zp *ec.Point

  ZBytes []byte
  StateIndex uint64
  StateBit uint32
}

func NewDualECDRBG(c *ec.PrimeCurve, qx, qy, seed *big.Int) *DualECDRBG {
  g := DualECDRBG{c, &ec.Point{qx, qy, true}, big.NewInt(0), nil, big.NewInt(0), nil, nil, 0, 0}
  g.seed(seed)
  return &g
}

// Seed is twice security_strength bits long (at least 256)
func (g *DualECDRBG) seed(seed *big.Int) {
  g.S = seed
  // To match OpenSSL, we perform two rounds before giving any output.
  g.generate_number()
  g.generate_number()
}

func (g *DualECDRBG) generate_number() {
  g.StateIndex++

  g.Sp = g.C.ScalarMultiply(g.S, g.C.G)
  g.S = g.Sp.X

  g.Zp = g.C.ScalarMultiply(g.S, g.Q)
  g.Z = g.Zp.X

  // We use .Bytes() so as to extract Z in big-endian format. This matches
  // OpenSSL behaviour. We drop the first 16 bits by initialising g.StateBit
  // as 16.
  g.ZBytes = g.Z.Bytes()
  g.StateBit = 16
}

func (g *DualECDRBG) Selfcheck() error {
  if !g.C.Satisfied(g.Sp) {
    return errors.New("S is not on curve.")
  }
  if !g.C.Satisfied(g.Zp) {
    return errors.New("Z is not on curve.")
  }

  if !g.Sp.Finite {
    return errors.New("S is infinite.")
  }
  if !g.Zp.Finite {
    return errors.New("Z is infinite.")
  }

  return nil
}

func (g *DualECDRBG) Bit() uint32 {
  if int(g.StateBit) >= len(g.ZBytes) * 8 {
    g.generate_number()
  }

  z_byte := g.ZBytes[g.StateBit / 8]
  z_bit_shift := 7 - (g.StateBit % 8)
  z_bit := 1 & (z_byte >> z_bit_shift)

  g.StateBit++

  return uint32(z_bit)
}

func (g *DualECDRBG) Urand32() uint32 {
  v := uint32(0)
  for i := 0; i < 32; i++ {
    v = (v<<1) + g.Bit()
  }
  return v
}

func (g *DualECDRBG) Byte() byte {
  if g.StateBit % 8 != 0 {
    g.StateBit = ((g.StateBit + 8) >> 3) << 3
  }

  if int(g.StateBit) >= len(g.ZBytes) * 8 {
    g.generate_number()
  }

  z_byte := g.ZBytes[g.StateBit / 8]
  g.StateBit += 8
  return z_byte
}

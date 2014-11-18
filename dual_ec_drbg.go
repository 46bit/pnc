package pinocchio

import (
  "fmt"
  "math/big"
  "github.com/46bit/pinocchio/ec"
)

type DualECDRBG struct {
  C *ec.PrimeCurve
  Q *ec.Point

  Z *big.Int
  z_bytes []byte
  S *big.Int
  StateIndex uint64
  StateBit uint32
}

func NewDualECDRBG(c *ec.PrimeCurve, qx, qy, seed *big.Int) *DualECDRBG {
  g := DualECDRBG{c, &ec.Point{qx, qy, true}, big.NewInt(0), nil, big.NewInt(0), 0, 0}
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

  s := g.C.ScalarMultiply(g.S, g.C.G)
  z := g.C.ScalarMultiply(s.X, g.Q)

  fmt.Printf("s is\nx = %X\ny = %X\n\n", s.X, s.Y)
  fmt.Printf("z is\nx = %X\ny = %X\n\n", z.X, z.Y)

  if !g.C.Satisfied(s) {
    fmt.Printf("s = g.C.G * %X not on curve\n", g.S)
  }
  if !g.C.Satisfied(z) {
    fmt.Printf("z = g.Q * %X not on curve\n", s.X)
  }

  if !s.Finite {
    fmt.Println("s infinite")
  }
  if !z.Finite {
    fmt.Println("z infinite")
  }

  g.S = s.X
  g.Z = z.X

  g.z_bytes = g.Z.Bytes()
  g.StateBit = 16
  fmt.Printf("%x\n", g.z_bytes)
}

func (g *DualECDRBG) Bit() uint32 {
  z_byte := g.z_bytes[g.StateBit / 8]
  z_bit_shift := 7 - (g.StateBit % 8)
  bit := 1 & (z_byte >> z_bit_shift)

  g.StateBit++
  if int(g.StateBit) >= g.Z.BitLen() { // @TODO: think if the bitlen-16 is off by 1
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

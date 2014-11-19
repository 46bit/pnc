package pnc

// http://en.wikipedia.org/wiki/Linear_congruential_generator#Parameters_in_common_use
// Numerical Recipes  232 1664525 1013904223
const (
  lcg_modulus = 1<<32
  lcg_multiplier = 1664525
  lcg_increment = 1013904223
)

type LCG struct {
  State uint32
  StateBit uint32
  StateIndex uint64
}

func NewLCG() LCG {
  l := LCG{}
  return l
}

func (l *LCG) Seed(seed uint32) {
  l.State = seed
  l.generate_number()
}

func (l *LCG) generate_number() {
  l.StateIndex++
  l.StateBit = 0
  l.State = uint32((lcg_multiplier * uint64(l.State) + lcg_increment) % lcg_modulus)
}

func (l *LCG) Bit() uint32 {
  bit_index_from_left := (31 - l.StateBit)
  bit := 1 & (l.State>>bit_index_from_left)
  l.StateBit++
  if l.StateBit == 32 {
    l.generate_number()
  }
  return bit
}

func (l *LCG) Urand32() uint32 {
  v := uint32(0)
  for i := 0; i < 32; i++ {
    v = (v<<1) + l.Bit()
  }
  return v
}

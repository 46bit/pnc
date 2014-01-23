package pinocchio

// http://en.wikipedia.org/wiki/Linear_congruential_generator#Parameters_in_common_use
// Numerical Recipes  232 1664525 1013904223
const (
  lcg_modulus = 1<<32
  lcg_multiplier = 1664525
  lcg_increment = 1013904223
)

type LCG struct {
  State uint32
  index int
}

func NewLCG() LCG {
  l := LCG{}
  return l
}

func (l *LCG) Seed(seed uint32) {
  l.index = 0
  l.State = seed
}

func (l *LCG) generate_number() {
  l.index++
  l.State = uint32((lcg_multiplier * uint64(l.State) + lcg_increment) % lcg_modulus)
}

func (l *LCG) Urand32() uint32 {
  l.generate_number()
  return l.State
}

/*
func main() {
  l := NewLCG()
  l.Seed(0)

  for i := 0; i < 500; i++ {
    fmt.Printf("%d: %d\n", l.index, l.Urand32())
  }
}
*/

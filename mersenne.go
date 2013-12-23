package main

// with the restriction that 2nw − r − 1 is a Mersenne prime
const (
  mersenne_twister_n = 624 // degree of recurrence
  mersenne_twister_m = 397 // middle word, or the number of parallel sequences, 1 ≤ m ≤ n
  mersenne_twister_r = 31 // separation point of one word, or the number of bits of the lower bitmask, 0 ≤ r ≤ w - 1
  mersenne_twister_a = 0x9908B0DF // coefficients of the rational normal form twist matrix
  mersenne_twister_b = 0x9D2C5680 // TGFSR(R) tempering bitmasks
  mersenne_twister_c = 0xEFC60000 // TGFSR(R) tempering bitmasks
  mersenne_twister_s = 7 // TGFSR(R) tempering bit shifts
  mersenne_twister_t = 15 // TGFSR(R) tempering bit shifts
  mersenne_twister_u = 11 // additional Mersenne Twister tempering bit shifts
  mersenne_twister_l = 18 // additional Mersenne Twister tempering bit shifts
)

type MersenneTwister struct {
  State [mersenne_twister_n]uint32
  index int
}

func NewMersenneTwister() MersenneTwister {
  m := MersenneTwister{}
  return m
}

func (m *MersenneTwister) Seed(seed uint32) {
  m.index = 0
  m.State[0] = seed
  for i := 1; i < mersenne_twister_n; i++ {
    t := m.State[i - 1] ^ (m.State[i - 1] >> 30)
    m.State[i] = 1812433253 * (t + uint32(i))
  }
}

func (m *MersenneTwister) SeedFromUrand32s(urand32s [mersenne_twister_n]uint32) {
  m.index = 0
  for i := 0; i < len(urand32s); i++ {
    m.State[i] = m.Urand32ToState(urand32s[i])
  }
}

func (m *MersenneTwister) generate_numbers() {
  for i := 0; i < mersenne_twister_n; i++ {
    y := (m.State[i] & 0x80000000) + (m.State[(i + 1) % mersenne_twister_n] & 0x7fffffff)
    m.State[i] = m.State[(i + mersenne_twister_m) % mersenne_twister_n] ^ (y >> 1)
    if y % 2 != 0 {
      m.State[i] = m.State[i] ^ mersenne_twister_a
    }
  }
}

func (m *MersenneTwister) Urand32() uint32 {
  if m.index == 0 {
    m.generate_numbers()
  }

  y := m.State[m.index]
  y = y ^ (y >> mersenne_twister_u)
  y = y ^ ((y << mersenne_twister_s) & mersenne_twister_b)
  y = y ^ ((y << mersenne_twister_t) & mersenne_twister_c)
  y = y ^ (y >> mersenne_twister_l)

  m.index = (m.index + 1) % mersenne_twister_n
  return y
}

func (m *MersenneTwister) coded_xor(urand32 uint32, changebit_index uint32, xorbit_index uint32) uint32 {
  changebit_shift_index := 31 - changebit_index
  xorbit_shift_index := 31 - xorbit_index

  changebit := (urand32 >> changebit_shift_index) & 1
  xorbit := (urand32 >> xorbit_shift_index) & 1
  result := changebit ^ xorbit

  if changebit == 1 && result == 0 {
    urand32 -= (1 << changebit_shift_index)
  } else if changebit == 0 && result == 1 {
    urand32 += (1 << changebit_shift_index)
  }

  return urand32
}

func (m *MersenneTwister) coded_and_xor(urand32 uint32, uand32 uint32, changebit_index uint32, xorbit_index uint32) uint32 {
  xorbit_shift_index := 31 - xorbit_index
  changebit_shift_index := 31 - changebit_index

  changebit := (urand32 >> changebit_shift_index) & 1
  uandbit := (uand32 >> changebit_shift_index) & 1
  xorbit := (urand32 >> xorbit_shift_index) & 1
  result := xorbit & uandbit

  result = changebit ^ result

  if changebit == 1 && result == 0 {
    urand32 -= (1 << changebit_shift_index)
  } else if changebit == 0 && result == 1 {
    urand32 += (1 << changebit_shift_index)
  }

  return urand32
}

func (m *MersenneTwister) Urand32ToState(urand32 uint32) uint32 {
  for i := uint32(18); i <= 31; i++ {
    urand32 = m.coded_xor(urand32, i, i - 18)
  }

  for i := uint32(16); true; i-- {
    urand32 = m.coded_and_xor(urand32, mersenne_twister_c, i, i + 15)
    if i == 0 { break }
  }

  for i := uint32(24); true; i-- {
    urand32 = m.coded_and_xor(urand32, mersenne_twister_b, i, i + 7)
    if i == 0 { break }
  }

  for i := uint32(11); i <= 31; i++ {
    urand32 = m.coded_xor(urand32, i, i - 11)
  }

  return urand32
}

/*
func main() {
  m := NewMersenneTwister()
  m.Seed(0)

  var urand32s [624]uint32
  for i := 0; i < len(urand32s); i++ {
    urand32s[i] = m.Urand32()
  }

  m2 := NewMersenneTwister()
  m2.SeedFromUrand32s(urand32s)

  for i := 0; i < 10000; i++ {
    fmt.Printf("%d = %d\n", m.Urand32(), m2.Urand32())
  }
}
*/

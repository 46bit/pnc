package pnc

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
)

const (
  p = 98
  L = 32
  d = 100 * p
  r = 5000 * p
  lfsr_width = d * L
)

type TauswortheGenerator struct {
  Registers [lfsr_width]uint32
  StateBit uint32
  StateIndex uint64
}

func NewTauswortheGenerator() TauswortheGenerator {
  t := TauswortheGenerator{}
  return t
}

func (t *TauswortheGenerator) Seed(seed uint32) {
  fmt.Println("TauswortheGenerator.Seed: note that seeding this takes ages due to slow mixing")
  for pi := 0; pi < lfsr_width; pi++ {
    if (seed >> (uint32(pi) % 32)) % 2 == 1 {
      t.Registers[pi] = 1
    } else {
      t.Registers[pi] = 0
    }
  }
  t.StateBit = 32
  for ri := 0; ri < r; ri++ {
    t.Urand32()
  }
}

func (t *TauswortheGenerator) generate_number() {
  t.StateIndex++
  t.StateBit = 0
  v := (t.Registers[lfsr_width - 33] + t.Registers[lfsr_width - 48]) % 2
  for lfsri := 1; lfsri < lfsr_width; lfsri++ {
    t.Registers[lfsri - 1] = t.Registers[lfsri]
  }
  t.Registers[lfsr_width - 1] = v
}

func (t *TauswortheGenerator) Bit() uint32 {
  if t.StateBit >= 32 {
    t.generate_number()
  }
  bit := uint32(t.Registers[t.StateBit * d])
  t.StateBit++
  return bit
}

func (t *TauswortheGenerator) Urand32() uint32 {
  v := uint32(0)
  for Li := 0; Li < L; Li++ {
    v = v<<1 + t.Bit()
  }
  return v
}

func (t *TauswortheGenerator) AsJSON() ([]byte, error) {
  t_json, err := json.Marshal(t)
  if err != nil {
    return []byte{}, err
  }
  return t_json, nil
}

func NewTauswortheGeneratorFromJSON(tg_json []byte) (TauswortheGenerator, error) {
  var t TauswortheGenerator
  err := json.Unmarshal(tg_json, &t)
  if err != nil {
    return t, err
  }
  return t, nil
}

func NewTauswortheGeneratorFromJSONFile(jpath string) (TauswortheGenerator, error) {
  tausworthe_generator_json, err := ioutil.ReadFile(jpath)
  if err != nil {
    return TauswortheGenerator{}, err
  }
  return NewTauswortheGeneratorFromJSON(tausworthe_generator_json)
}

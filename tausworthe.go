package pinocchio

import (
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
  current_index uint32
  Registers [lfsr_width]uint32
}

func NewTauswortheGenerator() TauswortheGenerator {
  t := TauswortheGenerator{}
  return t
}

func (t *TauswortheGenerator) Seed(seed uint32) {
  fmt.Println("TauswortheGenerator.Seed: note that seeding this takes ages due to slow mixing")
  t.current_index = 0
  for pi := 0; pi < lfsr_width; pi++ {
    if (seed >> (uint32(pi) % 32)) % 2 == 1 {
      t.Registers[pi] = 1
    } else {
      t.Registers[pi] = 0
    }
  }
  for ri := 0; ri < r; ri++ {
    t.Urand32()
  }
}

func (t *TauswortheGenerator) generate_number() {
  t.current_index = 0
  v := (t.Registers[lfsr_width - 33] + t.Registers[lfsr_width - 48]) % 2
  for lfsri := 1; lfsri < lfsr_width; lfsri++ {
    t.Registers[lfsri - 1] = t.Registers[lfsri]
  }
  t.Registers[lfsr_width - 1] = v
}

func (t *TauswortheGenerator) Bit() uint32 {
  bit := uint32(t.Registers[t.current_index * d])
  t.current_index++
  if t.current_index == 32 {
    t.generate_number()
  }
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

/*
func main() {
  debug := true

  t := NewTauswortheGenerator()
  t.Seed(987234789)

  // Print the generator post-seed as JSON.
  //t_json, _ := t.AsJSON()
  //fmt.Println(string(t_json))

  // Load the generator from a JSON file. Handy for getting values quickly with a
  // known seed, given how long the mixing takes...
  // t, _ := NewTauswortheGeneratorFromJSONFile(os.Args[1])

  for wi := 0; wi < 20; wi++ {
    // @TODO FIX: debug bias towards repeat values. Generating intermediate values is a
    // BAD temporary workaround.
    t.Urand32()
    w_value := t.Urand32()

    fmt.Printf("%dth value is %d\n", wi, w_value)
  }

  if debug { fmt.Println("EOF") }
}
*/

package main

import (
  "fmt"
  "github.com/46bit/pnc"
)

func main() {
  t := pnc.NewTauswortheGenerator()
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
}

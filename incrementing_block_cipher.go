package main

import (
  "os"
  "fmt"
//  "errors"
  "math/big"
)

type IncrementingBlockCipher struct {
  IV *big.Int
  Plainblock string
  Key *big.Int
  Keybits int
}

func NewIncrementingBlockCipher() IncrementingBlockCipher {
  c := IncrementingBlockCipher{}
  return c
}

func (c *IncrementingBlockCipher) SetKey(keybits int, key *big.Int) {
  c.Keybits = keybits
  c.Key = key
}

func (c *IncrementingBlockCipher) SetIV(iv *big.Int) {
  c.IV = iv
}

/*func (c *IncrementingBlockCipher) RandomizeIV() {
  c.IV = // ?
}*/

func (c *IncrementingBlockCipher) BytesBlockCount(bytes_length int) int {
  return (bytes_length * 8 + c.Keybits - 1) / c.Keybits
}

func (c *IncrementingBlockCipher) EncodeBytes(plain_bytes []byte) ([]byte, error) {
  plain_bytes_length := len(plain_bytes)
  last_plain_bytes := c.BigIntAsBytes(c.IV)
  block_count := c.BytesBlockCount(plain_bytes_length)
  crypt_str := ""

  for i := 0; i < block_count; i++ {
    block_offset_start := i * (c.Keybits + 7) / 8
    block_offset_end := (i + 1) * (c.Keybits + 7) / 8
    if block_offset_end > plain_bytes_length {
      block_offset_end = plain_bytes_length
    }

    block_plain_bytes := plain_bytes[block_offset_start:block_offset_end]
    block_crypt_bigint, err := c.EncodeBlock(c.BytesAsBigInt(block_plain_bytes), c.BytesAsBigInt(last_plain_bytes))
    if err != nil {
      return nil, err
    }

    block_crypt_bytes := c.BigIntAsBytes(block_crypt_bigint)
    block_crypt_str := string(block_crypt_bytes)
    crypt_str += block_crypt_str
    fmt.Printf("Block %d. f('%s' ^ '%s') = '0x%x'\n", i, string(block_plain_bytes), c.BytesAsBigInt(last_plain_bytes), block_crypt_bytes)

    last_plain_bytes = block_plain_bytes
  }
  return []byte(crypt_str), nil
}

func (c *IncrementingBlockCipher) DecodeBytes(crypt_bytes []byte) ([]byte, error) {
  crypt_bytes_length := len(crypt_bytes)
  last_plain_bytes := c.BigIntAsBytes(c.IV)
  block_count := c.BytesBlockCount(crypt_bytes_length)
  plain_str := ""

  for i := 0; i < block_count; i++ {
    block_offset_start := i * (c.Keybits + 7) / 8
    if block_offset_start >= crypt_bytes_length {
      break
    }
    block_offset_end := (i + 1) * (c.Keybits + 7) / 8
    if block_offset_end > crypt_bytes_length {
      block_offset_end = crypt_bytes_length
    }

    block_crypt_bytes := crypt_bytes[block_offset_start:block_offset_end]
    block_plain_bigint, err := c.DecodeBlock(c.BytesAsBigInt(block_crypt_bytes), c.BytesAsBigInt(last_plain_bytes))
    if err != nil {
      return nil, err
    }

    block_plain_bytes := c.BigIntAsBytes(block_plain_bigint)
    block_plain_str := string(block_plain_bytes)
    plain_str += block_plain_str
    fmt.Printf("Round %d. f^-1('0x%x') ^ '%s' = '%s'\n", i, c.BytesAsBigInt(block_crypt_bytes), c.BytesAsBigInt(last_plain_bytes), block_plain_str)

    last_plain_bytes = block_plain_bytes
  }
  return []byte(plain_str), nil
}

func (c *IncrementingBlockCipher) EncodeBlock(plain *big.Int, previous_plain *big.Int) (*big.Int, error) {
  plain.Xor(plain, previous_plain) // CBC
  return plain.Xor(plain, c.Key), nil
}

func (c *IncrementingBlockCipher) DecodeBlock(block *big.Int, previous_plain *big.Int) (*big.Int, error) {
  block.Xor(block, c.Key)
  block.Xor(block, previous_plain)
  return block, nil
}

func (c *IncrementingBlockCipher) BigIntAsBytes(block *big.Int) []byte {
  block2 := big.NewInt(0)
  block2.SetString(block.String(), 10)

  plain := make([]byte, (block.BitLen() + 7) / 8)
  for i := 1; i <= len(plain); i++ {
    b := big.NewInt(0)
    b.And(block2, big.NewInt(0xff))
    block2.Rsh(block2, 8)
    plain[len(plain) - i] = byte(b.Int64())
  }
  return plain
}

func (c *IncrementingBlockCipher) BytesAsBigInt(block []byte) *big.Int {
  bigblock := big.NewInt(0)
  for i := 0; i < len(block); i++ {
    plainbyteint := big.NewInt(int64(block[uint(i)]))
    bigblock.Add(bigblock, plainbyteint)
    bigblock.Lsh(bigblock, 8)
  }
  bigblock.Rsh(bigblock, 8)
  return bigblock
}

func main() {
  c := NewIncrementingBlockCipher()

  key := big.NewInt(0)
  key.SetString(os.Args[1], 10)
  c.SetKey(len(os.Args[1]) * 8, key)

  iv := []byte(os.Args[2])[0:(len(os.Args[1]) * 4 + 7) / 8]
  c.SetIV(c.BytesAsBigInt(iv))

  fmt.Printf("%d-bit Key '0x%x' with %d-bit IV '0x%x'.\n\n", c.Keybits, c.Key, len(os.Args[2]) * 8, c.IV)

  input := []byte(os.Args[3])

  fmt.Printf("input: %d '%s'\n\n", len(input), input)
  cryptext, err := c.EncodeBytes(input)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Printf("cryptext: %d '0x%x'\n\n", len(cryptext), c.BytesAsBigInt(cryptext))
  decodedtext, err := c.DecodeBytes(cryptext)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Printf("decodedtext: %d '%s'\n", len(decodedtext), decodedtext)
}

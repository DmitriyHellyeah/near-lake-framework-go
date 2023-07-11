package types

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type BigInt struct {
	big.Int
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	var bigInt string
	err := json.Unmarshal(p, &bigInt)
	if err != nil {
		return err
	}
	if bigInt == "" {
		return nil
	}
	z := new(big.Int)

	_, ok := z.SetString(bigInt, 10)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = *z
	return nil
}
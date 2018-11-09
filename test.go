package main

import (
	"encoding/hex"
	"fmt"

	amino "github.com/tendermint/go-amino"
)

type Test struct {
	A int64
}

func (Test) t() {

}

type TestI interface {
	t()
}

type TestSlice struct {
	T  []Test
	TI []TestI
}

func main1() {
	cdc := amino.NewCodec()

	cdc.RegisterInterface((*TestI)(nil), nil)
	cdc.RegisterConcrete(Test{}, "test", nil)

	bz, err := cdc.MarshalBinary(Test{1})
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(bz))

	bz, err = cdc.MarshalBinary(TestSlice{
		T: []Test{
			Test{
				A: 1,
			},
			Test{
				A: 2,
			},
		},
		TI: []TestI{
			Test{
				A: 3,
			},
			Test{
				A: 4,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(bz))
}

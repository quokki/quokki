package main

import (
	"encoding/hex"
	"fmt"

	amino "github.com/tendermint/go-amino"
)

type TestBin struct {
	A int64 `binary:"fixed64"`
}

type TestBin2 struct {
	A int64
}

func main3() {
	cdc := amino.NewCodec()

	bz, err := cdc.MarshalBinary(TestBin{1})
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(bz))

	bz, err = cdc.MarshalBinary(TestBin2{1})
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(bz))
}

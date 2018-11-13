package validator

import (
	codec "github.com/cosmos/cosmos-sdk/wire"
)

// Cosmos-sdk v0.25.0에서 wire를 codec으로 바꿨기 때문에 일단 codec으로 한다.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgValidator{}, "quokki/MsgValidator", nil)
}

var msgCdc = codec.NewCodec()

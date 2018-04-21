package answer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"
)

type AnswerMapper interface {
	GetAnswer(ctx sdk.Context, address sdk.Address) Answer
	CreateAnswer(ctx sdk.Context, question sdk.Address, writer sdk.Address, content string) sdk.Result
	UpdateAnswer(ctx sdk.Context, address sdk.Address, writer sdk.Address, content string) sdk.Result
}

type Answer interface {
	NewAddress(question sdk.Address, index int64)
	GetAddress() sdk.Address
	GetQuestion() sdk.Address
	GetIndex() int64
	GetWriter() sdk.Address
	GetCreateBlockHeight() int64
	GetContent(ctx sdk.Context, key sdk.StoreKey) string
	SetContent(ctx sdk.Context, key sdk.StoreKey, content string)
}

func RegisterWire() {
	const answerTypeBase = 0x1
	var _ = oldwire.RegisterInterface(
		struct{ Answer }{},
		oldwire.ConcreteType{&BaseAnswer{}, answerTypeBase},
	)
}

type AnswerDecoder func([]byte) (Answer, error)

func GetAnswerDecoder(cdc *wire.Codec) AnswerDecoder {
	return func(bytes []byte) (Answer, error) {
		var answer = &BaseAnswer{}

		err := cdc.UnmarshalBinary(bytes, &answer)
		return answer, err
	}
}

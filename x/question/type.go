package question

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"
)

type QuestionMapper interface {
	GetQuestion(ctx sdk.Context, address sdk.Address) Question
	GetQuestionsAt(ctx sdk.Context, blockHeight int64) []sdk.Address
	CreateQuestion(ctx sdk.Context, writer sdk.Address, partaker sdk.Address, title string, content string, language string, category string, tags []string) sdk.Result
	UpdateQuestion(ctx sdk.Context, writer sdk.Address, address sdk.Address, title string, content string, language string, category string, tags []string) sdk.Result

	GetAnswerNum(ctx sdk.Context, address sdk.Address) int64
	IncreaseAnswerNum(ctx sdk.Context, address sdk.Address)
}

type Question interface {
	NewAddress(ctx sdk.Context, title string, content string)
	GetAddress() sdk.Address
	GetWriter() sdk.Address
	GetPartaker() sdk.Address
	GetCreateBlockHeight() int64
	GetTitle(ctx sdk.Context, key sdk.StoreKey) string
	SetTitle(ctx sdk.Context, key sdk.StoreKey, title string)
	GetContent(ctx sdk.Context, key sdk.StoreKey) string
	SetContent(ctx sdk.Context, key sdk.StoreKey, content string)
	GetLanguage(ctx sdk.Context, key sdk.StoreKey) string
	SetLanguage(ctx sdk.Context, key sdk.StoreKey, language string)
	GetCategory(ctx sdk.Context, key sdk.StoreKey) string
	SetCategory(ctx sdk.Context, key sdk.StoreKey, category string)
	GetTags(ctx sdk.Context, key sdk.StoreKey, cdc *wire.Codec) []string
	SetTags(ctx sdk.Context, key sdk.StoreKey, cdc *wire.Codec, tags []string)
}

func RegisterWire() {
	const questTypeBase = 0x1
	var _ = oldwire.RegisterInterface(
		struct{ Question }{},
		oldwire.ConcreteType{&BaseQuestion{}, questTypeBase},
	)
}

type QuestionDecoder func([]byte) (Question, error)

func GetQuestionDecoder(cdc *wire.Codec) QuestionDecoder {
	return func(bytes []byte) (Question, error) {
		var question = &BaseQuestion{}

		err := cdc.UnmarshalBinary(bytes, &question)
		return question, err
	}
}

package answer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quokki/quokki/util"
)

var _ Answer = (*BaseAnswer)(nil)

type BaseAnswer struct {
	Address           sdk.Address `json:"address"`
	Question          sdk.Address `json:"question"`
	Index             int64       `json:"index"`
	Writer            sdk.Address `json:"writer"`
	CreateBlockHeight int64       `json:"createBlockHeight"`
}

func (answer *BaseAnswer) NewAddress(question sdk.Address, index int64) {
	answer.Question = question
	answer.Index = index

	answer.Address = util.GetAddressIndexHash(question, index, "answer")
}

func (answer BaseAnswer) GetAddress() sdk.Address {
	return answer.Address
}

func (answer BaseAnswer) GetQuestion() sdk.Address {
	return answer.Question
}

func (answer BaseAnswer) GetIndex() int64 {
	return answer.Index
}

func (answer BaseAnswer) GetWriter() sdk.Address {
	return answer.Writer
}

func (answer BaseAnswer) GetCreateBlockHeight() int64 {
	return answer.CreateBlockHeight
}

func (answer BaseAnswer) GetContent(ctx sdk.Context, key sdk.StoreKey) string {
	store := ctx.KVStore(key)
	b := store.Get(append([]byte("content-"), answer.GetAddress()...))
	return string(b)
}

func (answer BaseAnswer) SetContent(ctx sdk.Context, key sdk.StoreKey, content string) {
	store := ctx.KVStore(key)
	store.Set(append([]byte("content-"), answer.GetAddress()...), []byte(content))
}

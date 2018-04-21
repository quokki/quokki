package question

import (
	"encoding/binary"

	"golang.org/x/crypto/ripemd160"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var _ Question = (*BaseQuestion)(nil)

type BaseQuestion struct {
	Address           sdk.Address `json:"address"`
	Writer            sdk.Address `json:"writer"`
	Partaker          sdk.Address `json:"partaker"`
	CreateBlockHeight int64       `json:"createBlockHeight"`
}

func (question *BaseQuestion) NewAddress(ctx sdk.Context, title string, content string) {
	bytes := make([]byte, 40)
	binary.LittleEndian.PutUint64(bytes, uint64(ctx.BlockHeight()))
	bytes = append(bytes, ctx.ChainID()...)
	bytes = append(bytes, question.Writer...)
	bytes = append(bytes, question.Partaker...)
	bytes = append(bytes, title...)
	bytes = append(bytes, content...)
	h := ripemd160.New()
	h.Write(bytes)
	question.Address = h.Sum(nil)
}

func (question BaseQuestion) GetAddress() sdk.Address {
	return question.Address
}

func (question BaseQuestion) GetWriter() sdk.Address {
	return question.Writer
}

func (question BaseQuestion) GetPartaker() sdk.Address {
	return question.Partaker
}

func (question BaseQuestion) GetCreateBlockHeight() int64 {
	return question.CreateBlockHeight
}

func (question BaseQuestion) GetTitle(ctx sdk.Context, key sdk.StoreKey) string {
	store := ctx.KVStore(key)
	b := store.Get(append([]byte("title-"), question.Address...))
	return string(b)
}

func (question BaseQuestion) SetTitle(ctx sdk.Context, key sdk.StoreKey, title string) {
	store := ctx.KVStore(key)
	store.Set(append([]byte("title-"), question.Address...), []byte(title))
}

func (question BaseQuestion) GetContent(ctx sdk.Context, key sdk.StoreKey) string {
	store := ctx.KVStore(key)
	b := store.Get(append([]byte("content-"), question.Address...))
	return string(b)
}

func (question BaseQuestion) SetContent(ctx sdk.Context, key sdk.StoreKey, content string) {
	store := ctx.KVStore(key)
	store.Set(append([]byte("content-"), question.Address...), []byte(content))
}

func (question BaseQuestion) GetLanguage(ctx sdk.Context, key sdk.StoreKey) string {
	store := ctx.KVStore(key)
	b := store.Get(append([]byte("language-"), question.Address...))
	return string(b)
}

func (question BaseQuestion) SetLanguage(ctx sdk.Context, key sdk.StoreKey, language string) {
	store := ctx.KVStore(key)
	store.Set(append([]byte("language-"), question.Address...), []byte(language))
}

func (question BaseQuestion) GetCategory(ctx sdk.Context, key sdk.StoreKey) string {
	store := ctx.KVStore(key)
	b := store.Get(append([]byte("category-"), question.Address...))
	return string(b)
}

func (question BaseQuestion) SetCategory(ctx sdk.Context, key sdk.StoreKey, category string) {
	store := ctx.KVStore(key)
	store.Set(append([]byte("category-"), question.Address...), []byte(category))
}

func (question BaseQuestion) GetTags(ctx sdk.Context, key sdk.StoreKey, cdc *wire.Codec) (result []string) {
	store := ctx.KVStore(key)
	b := store.Get(append([]byte("tags-"), question.Address...))
	result = []string{}
	err := cdc.UnmarshalBinary(b, &result)
	if err != nil {
		result = []string{}
	}
	return
}

func (question BaseQuestion) SetTags(ctx sdk.Context, key sdk.StoreKey, cdc *wire.Codec, tags []string) {
	store := ctx.KVStore(key)
	b, err := cdc.MarshalBinary(tags)
	if err == nil {
		store.Set(append([]byte("tags-"), question.Address...), b)
	}
}

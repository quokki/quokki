package comment

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quokki/quokki/util"
)

var _ Comment = (*BaseComment)(nil)

type BaseComment struct {
	Address           sdk.Address `json:"address"`
	Target            sdk.Address `json:"target"`
	Index             int64       `json:"index"`
	Writer            sdk.Address `json:"writer"`
	Content           string      `json:"content"`
	CreateBlockHeight int64       `json:"createBlockHeight"`
}

func (comment *BaseComment) NewAddress(target sdk.Address, index int64) {
	comment.Target = target
	comment.Index = index

	comment.Address = util.GetAddressIndexHash(target, index, "comment")
}

func (comment BaseComment) GetAddress() sdk.Address {
	return comment.Address
}

func (comment BaseComment) GetTarget() sdk.Address {
	return comment.Target
}

func (comment BaseComment) GetIndex() int64 {
	return comment.Index
}

func (comment BaseComment) GetWriter() sdk.Address {
	return comment.Writer
}

func (comment BaseComment) GetContent() string {
	return comment.Content
}

func (comment BaseComment) GetCreateBlockHeight() int64 {
	return comment.CreateBlockHeight
}

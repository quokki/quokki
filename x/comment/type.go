package comment

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"
)

type CommentTypeToInfo map[string]CommentInfo

type CommentInfo struct {
	Key            sdk.StoreKey
	CollectionName string
}

type CommentMapper interface {
	GetComment(ctx sdk.Context, address sdk.Address) Comment
	CreateComment(ctx sdk.Context, _type string, target sdk.Address, writer sdk.Address, content string) sdk.Result

	GetCommentNum(ctx sdk.Context, target sdk.Address) int64
	IncreaseCommentNum(ctx sdk.Context, target sdk.Address)
}

type Comment interface {
	NewAddress(target sdk.Address, index int64)
	GetAddress() sdk.Address
	GetTarget() sdk.Address
	GetIndex() int64
	GetWriter() sdk.Address
	GetContent() string
	GetCreateBlockHeight() int64
}

func RegisterWire() {
	const commentTypeBase = 0x1
	var _ = oldwire.RegisterInterface(
		struct{ Comment }{},
		oldwire.ConcreteType{&BaseComment{}, commentTypeBase},
	)
}

type CommentDecoder func([]byte) (Comment, error)

func GetCommentDecoder(cdc *wire.Codec) CommentDecoder {
	return func(bytes []byte) (Comment, error) {
		var comment = &BaseComment{}

		err := cdc.UnmarshalBinary(bytes, &comment)
		return comment, err
	}
}

package comment

import (
	"bytes"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"

	"github.com/quokki/quokki/db"
)

/*
Comments have target address, and they also have their address.
TODO: Comments of comments.
*/

var _ CommentMapper = (*BaseCommentMapper)(nil)

type BaseCommentMapper struct {
	key               sdk.StoreKey
	cdc               *wire.Codec
	commentTypeToInfo CommentTypeToInfo
}

func NewCommentMapper(key sdk.StoreKey, commentTypeToInfo CommentTypeToInfo) CommentMapper {
	cdc := wire.NewCodec()
	return BaseCommentMapper{
		key:               key,
		cdc:               cdc,
		commentTypeToInfo: commentTypeToInfo,
	}
}

func (mapper BaseCommentMapper) GetComment(ctx sdk.Context, address sdk.Address) Comment {
	store := ctx.KVStore(mapper.key)
	bz := store.Get(address)
	if len(bz) == 0 {
		return nil
	}
	question := mapper.decodeComment(bz)
	return question
}

func (mapper BaseCommentMapper) CreateComment(ctx sdk.Context, _type string, target sdk.Address, writer sdk.Address, content string) sdk.Result {
	info, ok := mapper.commentTypeToInfo[_type]
	if ok == false {
		return sdk.ErrInternal("Invalid type").Result()
	}

	store := ctx.KVStore(info.Key)
	bz := store.Get(target)
	if len(bz) == 0 {
		return sdk.ErrInternal("Target does not exist").Result()
	}

	if ctx.IsCheckTx() {
		return sdk.Result{}
	}

	comment := BaseComment{}
	comment.Content = content
	comment.Writer = writer
	comment.CreateBlockHeight = ctx.BlockHeight()
	commentNum := mapper.GetCommentNum(ctx, target)
	if commentNum >= 1000 {
		return sdk.ErrInternal("This article already has too many comment").Result()
	}
	comment.NewAddress(target, commentNum)
	mapper.IncreaseCommentNum(ctx, target)
	commentNum++

	store = ctx.KVStore(mapper.key)
	bz = mapper.encodeComment(&comment)
	store.Set(comment.GetAddress(), bz)

	subData := make(map[string]interface{})
	db.Insert(ctx, "comments", comment, subData)
	db.UpdateSilently(ctx, info.CollectionName, map[string]interface{}{"address": target.String()}, map[string]interface{}{"commentNum": commentNum}, nil)

	return sdk.Result{Data: comment.GetAddress()}
}

func (mapper BaseCommentMapper) GetCommentNum(ctx sdk.Context, target sdk.Address) int64 {
	store := ctx.KVStore(mapper.key)
	b := store.Get(append([]byte("comment_num-"), target...))
	if len(b) == 0 || len(b) != 8 {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(b))
}

func (mapper BaseCommentMapper) IncreaseCommentNum(ctx sdk.Context, target sdk.Address) {
	store := ctx.KVStore(mapper.key)
	i := mapper.GetCommentNum(ctx, target)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i+1))
	store.Set(append([]byte("comment_num-"), target...), b)
}

func (mapper BaseCommentMapper) encodeComment(comment Comment) []byte {
	bz, err := mapper.cdc.MarshalBinary(comment)
	if err != nil {
		panic(err)
	}
	return bz
}

func (mapper BaseCommentMapper) decodeComment(bz []byte) Comment {
	r, n, err := bytes.NewBuffer(bz), new(int), new(error)
	commentI := oldwire.ReadBinary(struct{ Comment }{}, r, len(bz), n, err)
	if *err != nil {
		panic(*err)
	}

	comment := commentI.(struct{ Comment }).Comment
	return comment
}

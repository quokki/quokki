package article

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	codec "github.com/cosmos/cosmos-sdk/wire"
)

type Keeper struct {
	cdc       *codec.Codec
	storeKey  sdk.StoreKey
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  key,
		codespace: codespace,
	}
}

func (keeper Keeper) NewArticle(ctx sdk.Context, writer sdk.AccAddress, parent []byte, payload string) (sdk.Tags, sdk.Error) {
	tags := sdk.EmptyTags()

	article := Article{
		Id:        []byte{},
		Writer:    writer,
		Parent:    parent,
		Sequence:  0,
		CreatedAt: ctx.BlockHeader().Time,
		Payload:   payload,
	}

	err := keeper.assignArticleId(ctx, &article)
	if err != nil {
		return sdk.EmptyTags(), err
	}
	err = keeper.SetArticle(ctx, article)
	if err != nil {
		return sdk.EmptyTags(), err
	}

	tags.AppendTag("new_article", article.Id)
	return tags, nil
}

func (keeper Keeper) assignArticleId(ctx sdk.Context, article *Article) sdk.Error {
	if len(article.Id) > 0 {
		return ErrAssignedArticle(keeper.codespace, article.Id)
	}

	parent, err := keeper.GetArticle(ctx, article.Parent)
	if err != nil {
		return err
	}

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, parent.Sequence)

	article.Id = append(article.Parent, bz...)
	parent.Sequence++

	err = keeper.SetArticle(ctx, parent)
	if err != nil {
		return err
	}

	return nil
}

func (keeper Keeper) GetArticle(ctx sdk.Context, id []byte) (Article, sdk.Error) {
	article := Article{}
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(append([]byte("article"), id...))
	if len(bz) == 0 {
		return Article{}, ErrNonexistentArticle(keeper.codespace, id)
	}
	err := keeper.cdc.UnmarshalBinaryBare(bz, &article)
	if err != nil {
		return Article{}, ErrInvalidArticle(keeper.codespace, id)
	}
	return article, nil
}

func (keeper Keeper) SetArticle(ctx sdk.Context, article Article) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz, err := keeper.cdc.MarshalBinaryBare(article)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(append([]byte("article"), article.Id...), bz)
	return nil
}

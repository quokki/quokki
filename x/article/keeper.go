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

func (keeper Keeper) NewArticle(ctx sdk.Context, writer sdk.AccAddress, parent []byte, payload []byte) (sdk.Tags, sdk.Error) {
	tags := sdk.EmptyTags()
	article := NewArticle(writer, parent, payload)
	if len(article.Parent) > 0 {
		_, err := keeper.GetArticle(ctx, article.Parent)
		if err != nil {
			return tags, err
		}
	}

	id := keeper.newArticleId(ctx, parent, article.Writer)
	err := keeper.SetArticle(ctx, id, article)
	if err != nil {
		return tags, err
	}

	tags.AppendTag("new_article", id)
	return tags, nil
}

// TODO: articleaddr~~~
func (keeper Keeper) newArticleId(ctx sdk.Context, parent []byte, writer sdk.AccAddress) []byte {
	store := ctx.KVStore(keeper.storeKey)
	var sequence uint64 = 0
	if store.Has([]byte("sequence")) {
		bz := store.Get([]byte("sequence"))
		sequence = binary.BigEndian.Uint64(bz)
	}

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, sequence)

	sequence++
	sbz := make([]byte, 8)
	binary.BigEndian.PutUint64(sbz, sequence)
	store.Set([]byte("sequence"), sbz)

	return append(parent, bz...)
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

func (keeper Keeper) SetArticle(ctx sdk.Context, id []byte, article Article) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz, err := keeper.cdc.MarshalBinaryBare(article)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(append([]byte("article"), id...), bz)
	return nil
}

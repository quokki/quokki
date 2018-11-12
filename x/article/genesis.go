package article

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genArticle GenesisArticle) {
	article := Article{
		Id:        []byte{},
		Writer:    genArticle.Writer,
		Parent:    []byte{},
		Sequence:  0,
		CreatedAt: ctx.BlockHeader().Time,
		Payload:   genArticle.Payload,
	}

	k.SetArticle(ctx, article)
}

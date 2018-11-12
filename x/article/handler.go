package article

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgWrite:
			return handleMsgWrite(ctx, k, msg)
		default:
			errMsg := "Unrecognized article Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgWrite(ctx sdk.Context, k Keeper, msg MsgWrite) sdk.Result {
	tags, err := k.NewArticle(ctx, msg.Writer, msg.Parent, msg.Payload)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

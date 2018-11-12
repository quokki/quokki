package faucet

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgFaucet:
			return handleMsgFaucet(ctx, k, msg)
		default:
			errMsg := "Unrecognized article Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgFaucet(ctx sdk.Context, k Keeper, msg MsgFaucet) sdk.Result {
	tags, err := k.NewFaucet(ctx, msg.Address)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

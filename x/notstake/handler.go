package notstake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
)

func NewHandler(keeper NotstakeKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SetMsg:
			return handleSetMsg(ctx, keeper, msg)
		default:
			return sdk.ErrUnknownRequest("No match for message type.").Result()
		}
	}
}

func handleSetMsg(ctx sdk.Context, keeper NotstakeKeeper, msg SetMsg) sdk.Result {
	err := keeper.SetValInfo(ctx, msg.PubKey, msg.Power, msg.Weight)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	valSet := abci.Validator{
		PubKey: msg.PubKey.Bytes(),
		Power:  msg.Power,
	}

	return sdk.Result{
		Code:             sdk.CodeOK,
		ValidatorUpdates: abci.Validators{valSet},
	}
}

package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgValidator:
			return handleMsgValidator(ctx, k, msg)
		default:
			errMsg := "Unrecognized article Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgValidator(ctx sdk.Context, k Keeper, msg MsgValidator) sdk.Result {
	if !k.IsValidAdmin(ctx, msg.Admin) {
		return sdk.ErrInternal("Invalid admin").Result()
	}
	tags, err := k.SetValidator(ctx, msg.Validator, msg.PubKey, msg.Power)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func EndBlocker(ctx sdk.Context, k Keeper) (ValidatorUpdates []abci.Validator) {
	validator, err := k.GetValidator(ctx)
	if err == nil && len(validator.Address) > 0 {
		abciVal := abci.Validator{
			Address: validator.Address,
			PubKey:  tmtypes.TM2PB.PubKey(validator.PubKey),
			Power:   validator.VotingPower,
		}

		ValidatorUpdates = append(ValidatorUpdates, abciVal)
	}
	k.ClearValidator(ctx)

	return
}

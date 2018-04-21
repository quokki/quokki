package simpleregister

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(am sdk.AccountMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case RegisterMsg:
			return handleRegisterMsg(ctx, am, msg)
		default:
			return sdk.ErrUnknownRequest("No match for message type.").Result()
		}
	}
}

func handleRegisterMsg(ctx sdk.Context, am sdk.AccountMapper, msg RegisterMsg) sdk.Result {
	addr := msg.PubKey.Address()
	if len(addr) != 20 {
		return sdk.ErrInvalidAddress("Invalid address").Result()
	}

	account := am.GetAccount(ctx, addr)
	if account == nil {
		account = am.NewAccountWithAddress(ctx, addr)
		account.SetPubKey(msg.PubKey)
		am.SetAccount(ctx, account)
	}
	return sdk.Result{}
}

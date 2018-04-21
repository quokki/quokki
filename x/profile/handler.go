package profile

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(pm ProfileMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ProfileMsg:
			if ctx.IsCheckTx() {
				return sdk.Result{}
			}
			return pm.SetProfile(ctx, msg.Profile, msg.Address)
		default:
			errMsg := fmt.Sprintf("Unrecognized profile Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

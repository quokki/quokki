package power

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(pk PowerKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case PowerUpMsg:
			return pk.PowerUp(ctx, msg.Address, msg.Quokki)
		case PowerDownMsg:
			return pk.PowerDown(ctx, msg.Address, msg.QuokkiPower)
		case PowerUseMsg:
			err := pk.PowerUse(ctx, msg.Address, msg.QuokkiPower, msg.Term)
			if err != nil {
				return sdk.ErrInternal(err.Error()).Result()
			}
			return sdk.Result{}
		default:
			errMsg := fmt.Sprintf("Unrecognized power Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

package vote

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(vk VoteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case VoteUpMsg:
			return vk.VoteUp(ctx, msg.Address, msg.Answer, msg.QuokkiPower)
		default:
			errMsg := fmt.Sprintf("Unrecognized vote Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

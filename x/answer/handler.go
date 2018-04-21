package answer

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(am AnswerMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreateAnswerMsg:
			return am.CreateAnswer(ctx, msg.Question, msg.Writer, msg.Content)
		case UpdateAnswerMsg:
			return am.UpdateAnswer(ctx, msg.Address, msg.Writer, msg.Content)
		default:
			errMsg := fmt.Sprintf("Unrecognized answer Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

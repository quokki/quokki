package comment

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(cm CommentMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreateCommentMsg:
			return cm.CreateComment(ctx, msg.Type_, msg.Target, msg.Writer, msg.Content)
		default:
			errMsg := fmt.Sprintf("Unrecognized comment Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

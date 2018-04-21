package question

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(qm QuestionMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreateQuestionMsg:
			if ctx.IsCheckTx() {
				return sdk.Result{}
			}
			return qm.CreateQuestion(ctx, msg.Writer, msg.Writer, msg.Title, msg.Content, msg.Language, msg.Category, msg.Tags)
		case UpdateQuestionMsg:
			return qm.UpdateQuestion(ctx, msg.Writer, msg.Address, msg.Title, msg.Content, msg.Language, msg.Category, msg.Tags)
		default:
			errMsg := fmt.Sprintf("Unrecognized question Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

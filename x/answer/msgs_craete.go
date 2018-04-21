package answer

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MAX_CONTENT_LENGTH = 50000

type CreateAnswerMsg struct {
	Question sdk.Address `json:"question"`
	Writer   sdk.Address `json:"writer"`
	Content  string      `json:"content"`
}

var _ sdk.Msg = CreateAnswerMsg{}

func validateBasicAnswerMsg(content string) sdk.Error {
	err := ""
	if utf8.RuneCountInString(content) > MAX_CONTENT_LENGTH || len(content) > MAX_CONTENT_LENGTH*4 {
		err += fmt.Sprintf("Content must be shorter than %d characters.", MAX_CONTENT_LENGTH)
	}
	if err != "" {
		return sdk.NewError(401, err)
	}
	return nil
}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewCreateAnswerMsg(question sdk.Address, writer sdk.Address, content string) CreateAnswerMsg {
	return CreateAnswerMsg{Question: question, Writer: writer, Content: content}
}

// Implements Msg.
func (msg CreateAnswerMsg) Type() string { return "answer" }

// Implements Msg.
func (msg CreateAnswerMsg) ValidateBasic() sdk.Error {
	if len(msg.Writer) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	return validateBasicAnswerMsg(msg.Content)
}

func (msg CreateAnswerMsg) String() string {
	return fmt.Sprintf("CreateAnswerMsg{%v}", msg.Writer)
}

// Implements Msg.
func (msg CreateAnswerMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg CreateAnswerMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg CreateAnswerMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Writer}
}

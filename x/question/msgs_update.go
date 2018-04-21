package question

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UpdateQuestionMsg struct {
	Address  sdk.Address `json:"question"`
	Writer   sdk.Address `json:"writer"`
	Title    string      `json:"title"`
	Content  string      `json:"content"`
	Language string      `json:"language"`
	Category string      `json:"category"`
	Tags     []string    `json:"tags"`
}

var _ sdk.Msg = UpdateQuestionMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewUpdateQuestionMsg(address sdk.Address, writer sdk.Address, title string, content string, language string, category string, tags []string) UpdateQuestionMsg {
	return UpdateQuestionMsg{Address: address, Writer: writer, Title: title, Content: content, Language: language, Category: category, Tags: tags}
}

// Implements Msg.
func (msg UpdateQuestionMsg) Type() string { return "question" }

// Implements Msg.
func (msg UpdateQuestionMsg) ValidateBasic() sdk.Error {
	if len(msg.Writer) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	return validateBasicQuestionMsg(msg.Title, msg.Content, msg.Language, msg.Category, msg.Tags)
}

func (msg UpdateQuestionMsg) String() string {
	return fmt.Sprintf("UpdateQuestionMsg{%v}", msg.Writer)
}

// Implements Msg.
func (msg UpdateQuestionMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg UpdateQuestionMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg UpdateQuestionMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Writer}
}

package answer

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UpdateAnswerMsg struct {
	Address sdk.Address `json:"answer"`
	Writer  sdk.Address `json:"writer"`
	Content string      `json:"content"`
}

var _ sdk.Msg = UpdateAnswerMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewUpdateAnswerMsg(address sdk.Address, writer sdk.Address, content string) UpdateAnswerMsg {
	return UpdateAnswerMsg{Address: address, Writer: writer, Content: content}
}

// Implements Msg.
func (msg UpdateAnswerMsg) Type() string { return "answer" }

// Implements Msg.
func (msg UpdateAnswerMsg) ValidateBasic() sdk.Error {
	if len(msg.Writer) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	return validateBasicAnswerMsg(msg.Content)
}

func (msg UpdateAnswerMsg) String() string {
	return fmt.Sprintf("UpdateAnswerMsg{%v}", msg.Writer)
}

// Implements Msg.
func (msg UpdateAnswerMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg UpdateAnswerMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg UpdateAnswerMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Writer}
}

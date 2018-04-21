package comment

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MAX_CONTENT_LENGTH = 1000

type CreateCommentMsg struct {
	Type_   string      `json:"type"`
	Target  sdk.Address `json:"target"`
	Writer  sdk.Address `json:"writer"`
	Content string      `json:"content"`
}

var _ sdk.Msg = CreateCommentMsg{}

func NewCreateCommentMsg(_type string, target sdk.Address, writer sdk.Address, content string) CreateCommentMsg {
	return CreateCommentMsg{Type_: _type, Target: target, Writer: writer, Content: content}
}

// Implements Msg.
func (msg CreateCommentMsg) Type() string { return "comment" }

// Implements Msg.
func (msg CreateCommentMsg) ValidateBasic() sdk.Error {
	if len(msg.Type_) == 0 {
		return sdk.ErrInternal("Invalid type")
	}

	if len(msg.Target) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	if len(msg.Writer) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	err := ""
	if len(msg.Content) == 0 {
		err += fmt.Sprintf("Content should not be empty.")
	}
	if utf8.RuneCountInString(msg.Content) > MAX_CONTENT_LENGTH || len(msg.Content) > MAX_CONTENT_LENGTH*4 {
		err += fmt.Sprintf("Content must be shorter than %d characters.", MAX_CONTENT_LENGTH)
	}
	if err != "" {
		return sdk.NewError(401, err)
	}
	return nil
}

func (msg CreateCommentMsg) String() string {
	return fmt.Sprintf("CreateCommentMsg{%v}", msg.Writer)
}

// Implements Msg.
func (msg CreateCommentMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg CreateCommentMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg CreateCommentMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Writer}
}

package article

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgWrite struct {
	Writer  sdk.AccAddress `json:"writer"`
	Parent  []byte         `json:"parent"`
	Payload []byte         `json:"payload"`
}

var _ sdk.Msg = MsgWrite{}

func NewMsgWrite(writer sdk.AccAddress, parent []byte, payload []byte) MsgWrite {
	return MsgWrite{
		Writer:  writer,
		Parent:  parent,
		Payload: payload,
	}
}

func (msg MsgWrite) Type() string {
	return "article"
}

func (msg MsgWrite) ValidateBasic() sdk.Error {
	if len(msg.Payload) >= 5000 {
		// TODO: codespace 처리하는거 cosmos-sdk에서 바뀌면 적용하기
		return ErrTooBigPayload(DefaultCodespace)
	}
	return nil
}

func (msg MsgWrite) GetSignBytes() []byte {
	// It seemed that if bytes is empty slice, that will be replace with nil when amino encoded.
	// Thus, after tx passed through tx decoder,
	// that will be decoded as nil,
	// so this msg may be different with cli and chain deamon without matching empty slice and nil as empty slice.
	if msg.Parent == nil {
		msg.Parent = []byte{}
	}
	if msg.Payload == nil {
		msg.Payload = []byte{}
	}
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgWrite) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Writer}
}

package faucet

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgFaucet struct {
	Address sdk.AccAddress `json:"address"`
}

var _ sdk.Msg = MsgFaucet{}

func NewMsgFaucet(address sdk.AccAddress) MsgFaucet {
	return MsgFaucet{
		Address: address,
	}
}

func (msg MsgFaucet) Type() string {
	return "faucet"
}

func (msg MsgFaucet) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgFaucet) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgFaucet) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

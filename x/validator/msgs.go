package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type MsgValidator struct {
	Admin     sdk.AccAddress `json:"admin"`
	Validator sdk.ValAddress `json:"validator"`
	PubKey    crypto.PubKey  `json:"pubKey"`
	Power     int64          `json:"power"`
}

var _ sdk.Msg = MsgValidator{}

func NewMsgValidator(admin sdk.AccAddress, validator sdk.ValAddress, pubKey crypto.PubKey, power int64) MsgValidator {
	return MsgValidator{
		Admin:     admin,
		Validator: validator,
		PubKey:    pubKey,
		Power:     power,
	}
}

func (msg MsgValidator) Type() string {
	return "validator"
}

func (msg MsgValidator) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgValidator) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(struct {
		Admin   sdk.AccAddress `json:"admin"`
		Address sdk.ValAddress `json:"address"`
		PubKey  string         `json:"pubkey"`
		Power   int64          `json:"power"`
	}{
		Admin:   msg.Admin,
		Address: msg.Validator,
		PubKey:  sdk.MustBech32ifyValPub(msg.PubKey),
		Power:   msg.Power,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Admin}
}

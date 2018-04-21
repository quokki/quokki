package notstake

import (
	"bytes"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

type SetMsg struct {
	PubKey crypto.PubKey `json:"pub_key"`
	Power  int64         `json:"power"`
	Weight int64         `json:"weight"`
	Admin  sdk.Address   `json:"admin"`
}

var _ sdk.Msg = SetMsg{}

func NewSetMsg(pubKey crypto.PubKey, power int64, weight int64, admin sdk.Address) SetMsg {
	return SetMsg{
		PubKey: pubKey,
		Power:  power,
		Weight: weight,
		Admin:  admin,
	}
}

func (msg SetMsg) Type() string {
	return "notstake"
}

func (msg SetMsg) ValidateBasic() sdk.Error {
	if msg.PubKey.Empty() {
		return sdk.ErrInvalidPubKey("BondMsg.PubKey must not be empty")
	}

	if msg.Power < 0 {
		return sdk.ErrInternal("Power must be greater than or equal 0")
	}

	if msg.Weight <= 0 {
		return sdk.ErrInternal("Weight must be greater than 0")
	}

	if len(msg.Admin) != 20 {
		return sdk.ErrInternal("Invalid admin address")
	}

	for _, adminAddr := range admins {
		if bytes.Equal(adminAddr, msg.Admin) {
			return nil
		}
	}
	return sdk.ErrUnauthorized("Unauthorized admin address")
}

func (msg SetMsg) Get(key interface{}) interface{} {
	return nil
}

func (msg SetMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg SetMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Admin}
}

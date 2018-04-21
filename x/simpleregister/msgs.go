package simpleregister

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

type RegisterMsg struct {
	PubKey crypto.PubKey `json:"pub_key"`
}

var _ sdk.Msg = RegisterMsg{}

func NewRegisterMsg(pubKey crypto.PubKey) RegisterMsg {
	return RegisterMsg{PubKey: pubKey}
}

// Implements Msg.
func (msg RegisterMsg) Type() string { return "register" }

// Implements Msg.
func (msg RegisterMsg) ValidateBasic() sdk.Error {
	if msg.PubKey.Empty() {
		return sdk.ErrInvalidPubKey("Invalid pub key")
	}

	return nil
}

func (msg RegisterMsg) String() string {
	return fmt.Sprintf("RegisterMsg{%v}", msg.PubKey)
}

// Implements Msg.
func (msg RegisterMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg RegisterMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg RegisterMsg) GetSigners() []sdk.Address {
	return []sdk.Address{}
}

package profile

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProfileMsg struct {
	Address sdk.Address `json:"address"`
	Profile UserProfile `json:"profile"`
}

var _ sdk.Msg = ProfileMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewProfileMsg(address sdk.Address, profile UserProfile) ProfileMsg {
	return ProfileMsg{Address: address, Profile: profile}
}

// Implements Msg.
func (msg ProfileMsg) Type() string { return "profile" }

// Implements Msg.
func (msg ProfileMsg) ValidateBasic() sdk.Error {
	if len(msg.Address) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	return msg.Profile.ValidateBasic()
}

func (msg ProfileMsg) String() string {
	return fmt.Sprintf("ProfileMsg{%v}", msg.Address)
}

// Implements Msg.
func (msg ProfileMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg ProfileMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg ProfileMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

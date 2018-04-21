package power

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PowerDownMsg struct {
	Address     sdk.Address `json:"address"`
	QuokkiPower int64       `json:"quokki_power"`
}

var _ sdk.Msg = PowerDownMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewPowerDownMsg(address sdk.Address, quokkiPower int64) PowerDownMsg {
	return PowerDownMsg{Address: address, QuokkiPower: quokkiPower}
}

// Implements Msg.
func (msg PowerDownMsg) Type() string { return "power" }

// Implements Msg.
func (msg PowerDownMsg) ValidateBasic() sdk.Error {
	if len(msg.Address) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	if msg.QuokkiPower <= 0 {
		return sdk.ErrInsufficientCoins("Should not be zero")
	}

	return nil
}

func (msg PowerDownMsg) String() string {
	return fmt.Sprintf("PowerDownMsg{%v}", msg.Address)
}

// Implements Msg.
func (msg PowerDownMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg PowerDownMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg PowerDownMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

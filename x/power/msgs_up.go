package power

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PowerUpMsg struct {
	Address sdk.Address `json:"address"`
	Quokki  int64       `json:"quokki"`
}

var _ sdk.Msg = PowerUpMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewPowerUpMsg(address sdk.Address, quokki int64) PowerUpMsg {
	return PowerUpMsg{Address: address, Quokki: quokki}
}

// Implements Msg.
func (msg PowerUpMsg) Type() string { return "power" }

// Implements Msg.
func (msg PowerUpMsg) ValidateBasic() sdk.Error {
	if len(msg.Address) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	if msg.Quokki <= 0 {
		return sdk.ErrInsufficientCoins("Should not be zero")
	}

	return nil
}

func (msg PowerUpMsg) String() string {
	return fmt.Sprintf("PowerUpMsg{%v}", msg.Address)
}

// Implements Msg.
func (msg PowerUpMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg PowerUpMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg PowerUpMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

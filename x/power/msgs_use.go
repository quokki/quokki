package power

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PowerUseMsg struct {
	Address     sdk.Address `json:"address"`
	QuokkiPower int64       `json:"quokki"`
	Term        int64       `json:"restore_term"`
}

var _ sdk.Msg = PowerUseMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewPowerUseMsg(address sdk.Address, quokkiPower int64, term int64) PowerUseMsg {
	return PowerUseMsg{Address: address, QuokkiPower: quokkiPower, Term: term}
}

// Implements Msg.
func (msg PowerUseMsg) Type() string { return "power" }

// Implements Msg.
func (msg PowerUseMsg) ValidateBasic() sdk.Error {
	return sdk.ErrInternal("Test only")
	if len(msg.Address) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	if msg.QuokkiPower <= 0 {
		return sdk.ErrInsufficientCoins("Should not be zero")
	}

	return nil
}

func (msg PowerUseMsg) String() string {
	return fmt.Sprintf("PowerUseMsg{%v}", msg.Address)
}

// Implements Msg.
func (msg PowerUseMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg PowerUseMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg PowerUseMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

package vote

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteUpMsg struct {
	Address     sdk.Address `json:"address"`
	Answer      sdk.Address `json:"answer"`
	QuokkiPower int64       `json:"quokki_power"`
}

var _ sdk.Msg = VoteUpMsg{}

func NewVoteUpMsg(address sdk.Address, answer sdk.Address, quokkiPower int64) VoteUpMsg {
	return VoteUpMsg{Address: address, Answer: answer, QuokkiPower: quokkiPower}
}

// Implements Msg.
func (msg VoteUpMsg) Type() string { return "vote" }

// Implements Msg.
func (msg VoteUpMsg) ValidateBasic() sdk.Error {
	if len(msg.Address) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	if len(msg.Answer) != 20 {
		return sdk.ErrInvalidAddress("Invalid answer address")
	}

	if msg.QuokkiPower <= 0 {
		return sdk.ErrInsufficientCoins("Should not be zero")
	}

	return nil
}

func (msg VoteUpMsg) String() string {
	return fmt.Sprintf("VoteUpMsg{%v}", msg.Address)
}

// Implements Msg.
func (msg VoteUpMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg VoteUpMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg VoteUpMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

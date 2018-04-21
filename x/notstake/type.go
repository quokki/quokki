package notstake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

var admins []sdk.Address

type NotstakeAdminGenesis struct {
	Admins []sdk.Address `json:"admins"`
}

func SetAdmins(_admins []sdk.Address) {
	admins = _admins
}

type ValInfo struct {
	PubKey crypto.PubKey
	Power  int64
	Weight int64
}

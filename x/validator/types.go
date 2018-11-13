package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type Validator struct {
	Address     sdk.ValAddress `json:"address"`
	PubKey      crypto.PubKey  `json:"pubKey"`
	VotingPower int64          `json:"votingPower"`
}

func SetAdmin(ctx sdk.Context, k Keeper, address sdk.AccAddress) {
	k.SetAdmin(ctx, address)
}

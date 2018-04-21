package power

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RestoreInfo struct {
	Address     sdk.Address
	QuokkiPower int64
}

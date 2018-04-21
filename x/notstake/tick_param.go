package notstake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type notstakeTickParam struct {
	TotalNotstakeSupply sdk.Rat `json:"total_notstake_supply"`
	InflationRate       sdk.Rat `json:"inflation_rate"`
}

func (keeper NotstakeKeeper) GetNotstakeTickParam(ctx sdk.Context) (param notstakeTickParam) {
	store := ctx.KVStore(keeper.key)
	bz := store.Get([]byte("notstake-tick-param"))
	if bz == nil {
		panic("notstake tick param not init")
	}
	err := keeper.cdc.UnmarshalBinary(bz, &param)
	if err != nil {
		panic("notstake tick param unmarshal error")
	}
	return
}

func (keeper NotstakeKeeper) SetNotstakeTickParam(ctx sdk.Context, param notstakeTickParam) {
	store := ctx.KVStore(keeper.key)
	bz, err := keeper.cdc.MarshalBinary(param)
	if err != nil {
		panic(err)
	}
	store.Set([]byte("notstake-tick-param"), bz)
}

type NotstakeTickParamGenesis struct {
	Param notstakeTickParam `json:"notstake_tick_param"`
}

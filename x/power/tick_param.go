package power

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type powerTickParam struct {
	TotalPowerSupply sdk.Rat `json:"total_power_supply"`
	UnusedSupply     sdk.Rat `json:"unused_supply"`
	InflationRate    sdk.Rat `json:"inflation_rate"`
}

func (keeper PowerKeeper) GetPowerTickParam(ctx sdk.Context) (param powerTickParam) {
	store := ctx.KVStore(keeper.key)
	bz := store.Get([]byte("power-tick-param"))
	if bz == nil {
		panic("power tick param not init")
	}
	err := keeper.cdc.UnmarshalBinary(bz, &param)
	if err != nil {
		panic("power tick param unmarshal error")
	}
	return
}

func (keeper PowerKeeper) SetPowerTickParam(ctx sdk.Context, param powerTickParam) {
	store := ctx.KVStore(keeper.key)
	bz, err := keeper.cdc.MarshalBinary(param)
	if err != nil {
		panic(err)
	}
	store.Set([]byte("power-tick-param"), bz)
}

type PowerTickParamGenesis struct {
	Param powerTickParam `json:"power_tick_param"`
}

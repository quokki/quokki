package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type voteTickParam struct {
	TotalVoteSupply sdk.Rat `json:"total_vote_supply"`
	UnusedSupply    sdk.Rat `json:"unused_supply"`
	InflationRate   sdk.Rat `json:"inflation_rate"`
}

func (keeper VoteKeeper) GetVoteTickParam(ctx sdk.Context) (param voteTickParam) {
	store := ctx.KVStore(keeper.key)
	bz := store.Get([]byte("vote-tick-param"))
	if bz == nil {
		panic("vote tick param not init")
	}
	err := keeper.cdc.UnmarshalBinary(bz, &param)
	if err != nil {
		panic("vote tick param unmarshal error")
	}
	return
}

func (keeper VoteKeeper) SetVoteTickParam(ctx sdk.Context, param voteTickParam) {
	store := ctx.KVStore(keeper.key)
	bz, err := keeper.cdc.MarshalBinary(param)
	if err != nil {
		panic(err)
	}
	store.Set([]byte("vote-tick-param"), bz)
}

type VoteTickParamGenesis struct {
	Param voteTickParam `json:"vote_tick_param"`
}

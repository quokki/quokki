package power

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quokki/quokki/types"
)

var BlocksPerProvision int64 = 14400 //4hours
var BlocksPerYear int64 = 31536000

/*
	오류나면 그냥 넘어가는데
	오류나면 Used를 다 Available로 바꿔주던가 해야될듯
*/
func (keeper PowerKeeper) Tick(ctx sdk.Context) {
	store := ctx.KVStore(keeper.key)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(ctx.BlockHeight()))
	dest := append([]byte("restore-"), b...)
	bz := store.Get(dest)

	restoreInfos := []RestoreInfo{}
	if bz != nil && len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &restoreInfos)
		if err != nil {
			return
		}
	}

	for i := 0; i < len(restoreInfos); i++ {
		restoreInfo := restoreInfos[i]
		account := keeper.am.GetAccount(ctx, restoreInfo.Address)
		if account == nil {
			continue
		}
		_quokkiPower, err := account.Get("QuokkiPower")
		if err != nil {
			continue
		}
		accountQuokkiPower, ok := _quokkiPower.(types.QuokkiPower)
		if ok == false {
			continue
		}

		resultQuokkiPower := accountQuokkiPower
		if restoreInfo.QuokkiPower > resultQuokkiPower.Used {
			restoreInfo.QuokkiPower = resultQuokkiPower.Used
		}
		resultQuokkiPower.Used -= restoreInfo.QuokkiPower
		resultQuokkiPower.Available += restoreInfo.QuokkiPower

		err = account.Set("QuokkiPower", resultQuokkiPower)
		if err != nil {
			continue
		}

		keeper.am.SetAccount(ctx, account)
	}

	if ctx.BlockHeight() > 0 && (ctx.BlockHeight()%BlocksPerProvision) == 0 {
		keeper.tickInflation(ctx)
	}
}

func (keeper PowerKeeper) tickInflation(ctx sdk.Context) {
	store := ctx.KVStore(keeper.key)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(ctx.BlockHeight()))
	dest := append([]byte("inflation-"), b...)
	bz := store.Get(dest)

	restoreInfos := []RestoreInfo{}
	if bz != nil && len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &restoreInfos)
		if err != nil {
			return
		}
	}

	totalUsedQuokkiPower := sdk.ZeroRat
	for _, info := range restoreInfos {
		totalUsedQuokkiPower = totalUsedQuokkiPower.Add(sdk.NewRat(info.QuokkiPower))
	}

	param := keeper.GetPowerTickParam(ctx)
	totalPowerSupply := param.TotalPowerSupply
	inflationQuokki := sdk.NewRat(totalPowerSupply.Mul(keeper.nextInflation(param)).Evaluate() * BlocksPerProvision)
	totalPowerSupply = totalPowerSupply.Add(inflationQuokki)
	unusedSupply := param.UnusedSupply

	if len(restoreInfos) == 0 || totalUsedQuokkiPower.IsZero() {
		unusedSupply = unusedSupply.Add(inflationQuokki)
	} else {
		unusedToInflation := sdk.NewRat(unusedSupply.Quo(sdk.NewRat(10)).Evaluate())
		inflationQuokki = inflationQuokki.Add(unusedToInflation)
		unusedSupply = unusedSupply.Sub(unusedToInflation)

		for _, info := range restoreInfos {
			calcQuokki := inflationQuokki.Quo(totalUsedQuokkiPower).Mul(sdk.NewRat(info.QuokkiPower)).Evaluate()
			keeper.ck.AddCoins(ctx, info.Address, sdk.Coins{sdk.Coin{Denom: "quokki", Amount: calcQuokki}})
		}
	}

	param.TotalPowerSupply = totalPowerSupply
	param.UnusedSupply = unusedSupply

	keeper.SetPowerTickParam(ctx, param)
}

func (keeper PowerKeeper) nextInflation(param powerTickParam) sdk.Rat {
	return param.InflationRate.Quo(sdk.NewRat(BlocksPerYear))
}

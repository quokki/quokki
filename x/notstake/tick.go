package notstake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var BlocksPerProvision int64 = 43200 //12hours
var BlocksPerYear int64 = 31536000

func (keeper NotstakeKeeper) Tick(ctx sdk.Context) {
	if ctx.BlockHeight() > 0 && ctx.BlockHeight()%BlocksPerProvision == 0 {
		valInfos := keeper.GetValInfos(ctx)
		var sumWeight int64 = 0
		for _, valInfo := range valInfos {
			sumWeight += valInfo.Weight
		}

		param := keeper.GetNotstakeTickParam(ctx)
		totalNotstakeSupply := param.TotalNotstakeSupply
		inflationQuokki := sdk.NewRat(totalNotstakeSupply.Mul(keeper.nextInflation(param)).Evaluate() * BlocksPerProvision)
		totalNotstakeSupply = totalNotstakeSupply.Add(inflationQuokki)

		for _, valInfo := range valInfos {
			address := valInfo.PubKey.Address()
			calcQuokki := inflationQuokki.Quo(sdk.NewRat(sumWeight)).Mul(sdk.NewRat(valInfo.Weight)).Evaluate()
			keeper.ck.AddCoins(ctx, address, sdk.Coins{sdk.Coin{Denom: "quokki", Amount: calcQuokki}})
		}

		param.TotalNotstakeSupply = totalNotstakeSupply
		keeper.SetNotstakeTickParam(ctx, param)
	}
}

func (keeper NotstakeKeeper) nextInflation(param notstakeTickParam) sdk.Rat {
	return param.InflationRate.Quo(sdk.NewRat(BlocksPerYear))
}

package power

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestPowerInflation(t *testing.T) {
	BlocksPerProvision = 10
	BlocksPerYear = 1000
	ctx, pk, _, ck := createTestInput(t, false)
	testHeight := 1000

	originQuokkis := []int64{}
	for _, addr := range addrs {
		originQuokkis = append(originQuokkis, ck.GetCoins(ctx, addr, nil).AmountOf("quokki"))
	}

	genesisParam := pk.GetPowerTickParam(ctx)
	assert.Equal(t, defaultParams(), genesisParam)

	for height := 0; height <= testHeight; height += 1 {
		ctx = ctx.WithBlockHeight(int64(height))
		pk.Tick(ctx)

		if int64(height)%BlocksPerProvision == 0 {
			for _, addr := range addrs {
				pk.PowerUse(ctx, addr, 10, BlocksPerProvision/2)
			}
		}
	}

	param := pk.GetPowerTickParam(ctx)
	supply := param.TotalPowerSupply.Evaluate()
	expectedSupply := genesisParam.TotalPowerSupply.Evaluate()
	expectedAddSupply := (sdk.NewRat(expectedSupply / BlocksPerYear).Mul(genesisParam.InflationRate)).Evaluate() * int64(testHeight)
	expectedSupply += expectedAddSupply
	assert.InDelta(t, expectedSupply, supply, float64(expectedSupply/100))

	quokkis := []int64{}
	for _, addr := range addrs {
		quokkis = append(quokkis, ck.GetCoins(ctx, addr, nil).AmountOf("quokki"))
	}

	addSupply := int64(0)
	for i, _ := range addrs {
		addSupply += quokkis[i] - originQuokkis[i]
	}
	assert.NotEqual(t, int64(0), addSupply)
	assert.InDelta(t, expectedAddSupply, addSupply+param.UnusedSupply.Evaluate(), float64(expectedAddSupply/10))
}

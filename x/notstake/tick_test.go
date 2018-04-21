package notstake

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNotstakeInflation(t *testing.T) {
	BlocksPerProvision = 10
	BlocksPerYear = 1000
	ctx, nk, ck := createTestInput(t, false)
	testHeight := 1000

	genesisParam := nk.GetNotstakeTickParam(ctx)
	assert.Equal(t, defaultParams(), genesisParam)

	for height := 1; height <= testHeight; height += 1 {
		ctx = ctx.WithBlockHeight(int64(height))
		nk.Tick(ctx)
	}

	param := nk.GetNotstakeTickParam(ctx)
	supply := param.TotalNotstakeSupply.Evaluate()
	expectedSupply := genesisParam.TotalNotstakeSupply.Evaluate()
	expectedAddSupply := (sdk.NewRat(expectedSupply / BlocksPerYear).Mul(genesisParam.InflationRate)).Evaluate() * int64(testHeight)
	expectedSupply += expectedAddSupply
	assert.InDelta(t, expectedSupply, supply, float64(expectedSupply/100))

	var sum int64 = 0
	var sumI int64 = 0
	for i, pk := range pks {
		sum += ck.GetCoins(ctx, pk.Address(), nil).AmountOf("quokki")
		sumI += int64(i + 1)
	}
	for i, pk := range pks {
		expected := sum / sumI * int64(i+1)
		assert.InDelta(t, expected, ck.GetCoins(ctx, pk.Address(), nil).AmountOf("quokki"), float64(expected/100))
	}

	assert.InDelta(t, expectedAddSupply, sum, float64(expectedAddSupply/10))
}

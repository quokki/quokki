package vote

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestVoteInflation(t *testing.T) {
	BlocksPerCompare = 10 //for testing
	BlocksCompareNum = 42
	BlocksPerYear = 2000
	testHeight := 1000
	ctx, vk, _, am, qm, ck := createTestInput(t, false)
	questions := []sdk.Address{}
	answers := []sdk.Address{}

	originQuokkis := []int64{}
	for _, addr := range addrs {
		originQuokkis = append(originQuokkis, ck.GetCoins(ctx, addr, nil).AmountOf("quokki"))
	}

	for height := 1; height <= testHeight; height += (rand.Intn(9) + 1) {
		ctx = ctx.WithBlockHeight(int64(height))
		for rand.Intn(100) < 50 {
			i := rand.Intn(len(addrs))
			result := qm.CreateQuestion(ctx, addrs[i], addrs[i], "", "", "", "", []string{})
			if result.IsOK() {
				questions = append(questions, result.Data)

				for rand.Intn(100) < 50 {
					i := rand.Intn(len(addrs))
					result := am.CreateAnswer(ctx, result.Data, addrs[i], "")

					if result.IsOK() {
						answers = append(answers, result.Data)

						for rand.Intn(100) < 50 {
							i := rand.Intn(len(addrs))
							vk.VoteUp(ctx, addrs[i], result.Data, int64(rand.Intn(1000)))
						}
					}
				}
			}
		}
	}

	genesisParam := vk.GetVoteTickParam(ctx)
	assert.Equal(t, defaultParams(), genesisParam)

	for height := 1; height <= testHeight; height += 1 {
		ctx = ctx.WithBlockHeight(int64(height))
		vk.Tick(ctx)
	}

	param := vk.GetVoteTickParam(ctx)
	supply := param.TotalVoteSupply.Evaluate()
	expectedSupply := genesisParam.TotalVoteSupply.Evaluate()
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
	assert.InDelta(t, expectedAddSupply, addSupply+param.UnusedSupply.Evaluate(), float64(expectedAddSupply/10))
}

package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quokki/quokki/db"
	"github.com/quokki/quokki/util"
)

var BlocksPerCompare int64 = 14400 //4hours
var BlocksCompareNum int64 = 42    //7days
var BlocksPerYear int64 = 31536000
var AnswerWeight sdk.Rat = sdk.NewRat(5, 100)

func (keeper VoteKeeper) Tick(ctx sdk.Context) {
	questions := []sdk.Address{}
	totalTickVote := sdk.NewRat(0)
	for blockHeight := ctx.BlockHeight() - BlocksPerCompare; blockHeight > 0 && blockHeight >= ctx.BlockHeight()-(BlocksPerCompare*BlocksCompareNum); blockHeight -= BlocksPerCompare {
		_questions := keeper.qm.GetQuestionsAt(ctx, blockHeight)
		for _, questAddr := range _questions {
			totalTickVote = totalTickVote.Add(sdk.NewRat(keeper.GetQuestionTotalVote(ctx, questAddr)))
		}
		questions = append(questions, _questions...)
	}
	totalTickVote = totalTickVote.Mul(sdk.OneRat.Add(AnswerWeight))

	param := keeper.GetVoteTickParam(ctx)
	totalVoteSupply := param.TotalVoteSupply
	inflationQuokki := sdk.NewRat(totalVoteSupply.Mul(keeper.nextInflation(param)).Evaluate())
	totalVoteSupply = totalVoteSupply.Add(inflationQuokki)
	unusedSupply := param.UnusedSupply

	if len(questions) == 0 || totalTickVote.IsZero() {
		unusedSupply = unusedSupply.Add(inflationQuokki)
	} else {
		unusedToInflation := sdk.NewRat(unusedSupply.Quo(sdk.NewRat(500)).Evaluate())
		inflationQuokki = inflationQuokki.Add(unusedToInflation)
		unusedSupply = unusedSupply.Sub(unusedToInflation)

		for _, questAddr := range questions {
			numAnswer := keeper.qm.GetAnswerNum(ctx, questAddr)
			var i int64 = 0
			for i = 0; i < numAnswer; i++ {
				answerAddr := util.GetAddressIndexHash(questAddr, i, "answer")
				answer := keeper.am.GetAnswer(ctx, answerAddr)
				vote := sdk.NewRat(keeper.GetAnswerVote(ctx, answerAddr))
				calcQuokki := inflationQuokki.Mul(vote.Quo(totalTickVote)).Evaluate()
				if answer == nil {
					unusedSupply = unusedSupply.Add(sdk.NewRat(calcQuokki))
					continue
				}
				writer := answer.GetWriter()
				if len(writer) == 20 {
					_, err := keeper.ck.AddCoins(ctx, writer, sdk.Coins{sdk.Coin{Amount: calcQuokki, Denom: "quokki"}})
					if err != nil {
						unusedSupply = unusedSupply.Add(sdk.NewRat(calcQuokki))
						continue
					}
					if calcQuokki > 0 {
						db.UpdateSilently(ctx, "answers",
							map[string]interface{}{"address": answerAddr.String()},
							map[string]interface{}{"$inc": map[string]interface{}{"earning": calcQuokki, "earningCount": 1}}, nil)
					}
				} else {
					unusedSupply = unusedSupply.Add(sdk.NewRat(calcQuokki))
					continue
				}
			}

			totalVote := sdk.NewRat(keeper.GetQuestionTotalVote(ctx, questAddr))
			calcQuokki := inflationQuokki.Mul(totalVote.Mul(AnswerWeight).Quo(totalTickVote)).Evaluate()
			question := keeper.qm.GetQuestion(ctx, questAddr)
			if question == nil {
				unusedSupply = unusedSupply.Add(sdk.NewRat(calcQuokki))
				continue
			}
			writer := question.GetWriter()
			if len(writer) == 20 {
				_, err := keeper.ck.AddCoins(ctx, writer, sdk.Coins{sdk.Coin{Amount: calcQuokki, Denom: "quokki"}})
				if err != nil {
					unusedSupply = unusedSupply.Add(sdk.NewRat(calcQuokki))
					continue
				}
				if calcQuokki > 0 {
					db.UpdateSilently(ctx, "questions",
						map[string]interface{}{"address": questAddr.String()},
						map[string]interface{}{"$inc": map[string]interface{}{"earning": calcQuokki, "earningCount": 1}}, nil)
				}
			} else {
				unusedSupply = unusedSupply.Add(sdk.NewRat(calcQuokki))
				continue
			}
		}
	}

	param.TotalVoteSupply = totalVoteSupply
	param.UnusedSupply = unusedSupply

	keeper.SetVoteTickParam(ctx, param)
}

func (keeper VoteKeeper) nextInflation(param voteTickParam) sdk.Rat {
	return param.InflationRate.Quo(sdk.NewRat(BlocksPerYear))
}

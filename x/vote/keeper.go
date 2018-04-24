package vote

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/quokki/quokki/db"
	"github.com/quokki/quokki/x/answer"
	"github.com/quokki/quokki/x/power"
	"github.com/quokki/quokki/x/question"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const VOTE_UP_RESTORE_TERM = 64800 //18 hours

type VoteKeeper struct {
	key sdk.StoreKey
	cdc *wire.Codec
	pk  power.PowerKeeper
	am  answer.AnswerMapper
	qm  question.QuestionMapper
	ck  bank.CoinKeeper
}

func NewVoteKeeper(key sdk.StoreKey, cdc *wire.Codec, pk power.PowerKeeper, am answer.AnswerMapper, qm question.QuestionMapper, ck bank.CoinKeeper) VoteKeeper {
	return VoteKeeper{key: key, cdc: cdc, pk: pk, am: am, qm: qm, ck: ck}
}

func (keeper VoteKeeper) VoteUp(ctx sdk.Context, address sdk.Address, answerAddress sdk.Address, quokkiPower int64) sdk.Result {
	if quokkiPower <= 0 {
		return sdk.ErrInvalidCoins("Should not be zero").Result()
	}

	answer := keeper.am.GetAnswer(ctx, answerAddress)
	if answer == nil {
		return sdk.ErrInvalidAddress("Answer does not exist").Result()
	}

	question := keeper.qm.GetQuestion(ctx, answer.GetQuestion())
	if question == nil {
		return sdk.ErrInvalidAddress("Answer does not exist").Result()
	}

	err := keeper.pk.PowerUse(ctx, address, quokkiPower, VOTE_UP_RESTORE_TERM)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	if ctx.IsCheckTx() {
		return sdk.Result{}
	}

	store := ctx.KVStore(keeper.key)
	var questionTotalVote int64 = 0
	bz := store.Get(append([]byte("question-total-vote"), question.GetAddress()...))
	if len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &questionTotalVote)
		if err != nil {
			return sdk.ErrInternal(err.Error()).Result()
		}
	}
	questionTotalVote += quokkiPower
	bz, err = keeper.cdc.MarshalBinary(questionTotalVote)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	store.Set(append([]byte("question-total-vote"), question.GetAddress()...), bz)

	var answerVote int64 = 0
	bz = store.Get(append([]byte("answer-vote"), answer.GetAddress()...))
	if len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &answerVote)
		if err != nil {
			return sdk.ErrInternal(err.Error()).Result()
		}
	}
	answerVote += quokkiPower
	bz, err = keeper.cdc.MarshalBinary(answerVote)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	store.Set(append([]byte("answer-vote"), answer.GetAddress()...), bz)

	db.UpdateSilently(ctx, "questions",
		map[string]interface{}{"address": question.GetAddress().String()},
		map[string]interface{}{"totalVote": questionTotalVote}, nil)
	db.UpdateSilently(ctx, "answers",
		map[string]interface{}{"address": answer.GetAddress().String()},
		map[string]interface{}{"vote": answerVote, "$inc": map[string]interface{}{"voteCount": 1}}, nil)

	db.Insert(ctx, "votes",
		map[string]interface{}{"address": answer.GetAddress().String(), "voter": address.String(), "vote": quokkiPower}, nil)

	return sdk.Result{}
}

func (keeper VoteKeeper) GetQuestionTotalVote(ctx sdk.Context, address sdk.Address) int64 {
	store := ctx.KVStore(keeper.key)
	var questionTotalVote int64 = 0
	bz := store.Get(append([]byte("question-total-vote"), address...))
	if len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &questionTotalVote)
		if err != nil {
			questionTotalVote = 0
		}
	}
	return questionTotalVote
}

func (keeper VoteKeeper) GetAnswerVote(ctx sdk.Context, address sdk.Address) int64 {
	store := ctx.KVStore(keeper.key)
	var answerVote int64 = 0
	bz := store.Get(append([]byte("answer-vote"), address...))
	if len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &answerVote)
		if err != nil {
			answerVote = 0
		}
	}
	return answerVote
}

package answer

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"

	"github.com/quokki/quokki/db"

	"github.com/quokki/quokki/x/question"
)

var _ AnswerMapper = (*BaseAnswerMapper)(nil)

type BaseAnswerMapper struct {
	key            sdk.StoreKey
	questionMapper question.QuestionMapper
	cdc            *wire.Codec
}

func NewAnswerMapper(key sdk.StoreKey, questionMapper question.QuestionMapper) AnswerMapper {
	cdc := wire.NewCodec()
	return BaseAnswerMapper{
		key:            key,
		questionMapper: questionMapper,
		cdc:            cdc,
	}
}

func (mapper BaseAnswerMapper) GetAnswer(ctx sdk.Context, address sdk.Address) Answer {
	store := ctx.KVStore(mapper.key)
	bz := store.Get(address)
	if len(bz) == 0 {
		return nil
	}
	question := mapper.decodeAnswer(bz)
	return question
}

func (mapper BaseAnswerMapper) CreateAnswer(ctx sdk.Context, question sdk.Address, writer sdk.Address, content string) sdk.Result {
	q := mapper.questionMapper.GetQuestion(ctx, question)
	if q == nil {
		return sdk.ErrInvalidAddress("Question does not exist").Result()
	}
	if ctx.IsCheckTx() {
		return sdk.Result{}
	}

	answer := BaseAnswer{}
	answer.Writer = writer
	answer.CreateBlockHeight = ctx.BlockHeight()
	answerNum := mapper.questionMapper.GetAnswerNum(ctx, question)
	if answerNum >= 100 {
		return sdk.ErrInternal("Question already has too many answer").Result()
	}
	answer.NewAddress(question, answerNum)
	mapper.questionMapper.IncreaseAnswerNum(ctx, question)
	answerNum++

	store := ctx.KVStore(mapper.key)
	bz := mapper.encodeAnswer(&answer)
	store.Set(answer.GetAddress(), bz)
	answer.SetContent(ctx, mapper.key, content)

	subData := make(map[string]interface{})
	subData["content"] = content
	db.Insert(ctx, "answers", answer, subData)
	db.UpdateSilently(ctx, "questions", map[string]interface{}{"address": question.String()}, map[string]interface{}{"answerNum": answerNum}, nil)

	return sdk.Result{Data: answer.GetAddress()}
}

func (mapper BaseAnswerMapper) UpdateAnswer(ctx sdk.Context, address sdk.Address, writer sdk.Address, content string) sdk.Result {
	store := ctx.KVStore(mapper.key)
	bz := store.Get(address)
	if len(bz) == 0 {
		return sdk.ErrInvalidAddress("Unrecognized answer address").Result()
	}
	answer := mapper.decodeAnswer(bz)
	if bytes.Equal(answer.GetWriter(), writer) == false {
		return sdk.ErrUnauthorized("Did not match writer").Result()
	}

	answer.SetContent(ctx, mapper.key, content)

	db.Update(ctx, "answers", map[string]interface{}{"address": answer.GetAddress().String()},
		map[string]interface{}{"content": content}, nil)
	return sdk.Result{}
}

func (mapper BaseAnswerMapper) encodeAnswer(answer Answer) []byte {
	bz, err := mapper.cdc.MarshalBinary(answer)
	if err != nil {
		panic(err)
	}
	return bz
}

func (mapper BaseAnswerMapper) decodeAnswer(bz []byte) Answer {
	r, n, err := bytes.NewBuffer(bz), new(int), new(error)
	answerI := oldwire.ReadBinary(struct{ Answer }{}, r, len(bz), n, err)
	if *err != nil {
		panic(*err)
	}

	answer := answerI.(struct{ Answer }).Answer
	return answer
}

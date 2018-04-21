package question

import (
	"bytes"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"

	"github.com/quokki/quokki/db"
)

var _ QuestionMapper = (*BaseQuestionMapper)(nil)

type BaseQuestionMapper struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewQuestionMapper(key sdk.StoreKey) QuestionMapper {
	cdc := wire.NewCodec()
	return BaseQuestionMapper{
		key: key,
		cdc: cdc,
	}
}

func (mapper BaseQuestionMapper) GetQuestion(ctx sdk.Context, address sdk.Address) Question {
	store := ctx.KVStore(mapper.key)
	bz := store.Get(address)
	if bz == nil {
		return nil
	}
	question := mapper.decodeQuestion(bz)
	return question
}

func (mapper BaseQuestionMapper) GetQuestionsAt(ctx sdk.Context, blockHeight int64) []sdk.Address {
	store := ctx.KVStore(mapper.key)
	questions := []sdk.Address{}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(blockHeight))
	bz := store.Get(append([]byte("question-at"), b...))
	if len(bz) > 0 {
		err := mapper.cdc.UnmarshalBinary(bz, &questions)
		if err != nil {
			questions = []sdk.Address{}
		}
	}
	return questions
}

func (mapper BaseQuestionMapper) CreateQuestion(ctx sdk.Context, writer sdk.Address, partaker sdk.Address, title string, content string, language string, category string, tags []string) sdk.Result {
	question := BaseQuestion{}
	question.Writer = writer
	question.Partaker = partaker
	question.CreateBlockHeight = ctx.BlockHeight()
	question.NewAddress(ctx, title, content)

	store := ctx.KVStore(mapper.key)

	if len(store.Get(question.GetAddress())) > 0 {
		return sdk.ErrInternal("Question address conflict").Result()
	}

	bz := mapper.encodeQuestion(&question)
	store.Set(question.GetAddress(), bz)
	question.SetTitle(ctx, mapper.key, title)
	question.SetContent(ctx, mapper.key, content)
	question.SetLanguage(ctx, mapper.key, language)
	question.SetCategory(ctx, mapper.key, category)
	question.SetTags(ctx, mapper.key, mapper.cdc, tags)

	questions := []sdk.Address{}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(ctx.BlockHeight()))
	bz = store.Get(append([]byte("question-at"), b...))
	if len(bz) > 0 {
		err := mapper.cdc.UnmarshalBinary(bz, &questions)
		if err != nil {
			questions = []sdk.Address{}
		}
	}
	questions = append(questions, question.GetAddress())
	bz, err := mapper.cdc.MarshalBinary(questions)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	store.Set(append([]byte("question-at"), b...), bz)

	subData := make(map[string]interface{})
	subData["title"] = title
	subData["content"] = content
	subData["language"] = language
	subData["category"] = category
	subData["tags"] = tags
	subData["answerNum"] = 0
	db.Insert(ctx, "questions", question, subData)

	return sdk.Result{Data: question.GetAddress()}
}

func (mapper BaseQuestionMapper) UpdateQuestion(ctx sdk.Context, writer sdk.Address, address sdk.Address, title string, content string, language string, category string, tags []string) sdk.Result {
	store := ctx.KVStore(mapper.key)
	bz := store.Get(address)
	if len(bz) == 0 {
		return sdk.ErrInvalidAddress("Unrecognized question address").Result()
	}
	question := mapper.decodeQuestion(bz)
	if bytes.Equal(question.GetWriter(), writer) == false {
		return sdk.ErrUnauthorized("Did not match writer").Result()
	}

	question.SetTitle(ctx, mapper.key, title)
	question.SetContent(ctx, mapper.key, content)
	question.SetLanguage(ctx, mapper.key, language)
	question.SetCategory(ctx, mapper.key, category)
	question.SetTags(ctx, mapper.key, mapper.cdc, tags)

	subData := make(map[string]interface{})
	subData["title"] = title
	subData["content"] = content
	subData["language"] = language
	subData["category"] = category
	subData["tags"] = tags
	db.Update(ctx, "questions", map[string]interface{}{"address": question.GetAddress().String()}, subData, nil)

	return sdk.Result{}
}

func (mapper BaseQuestionMapper) GetAnswerNum(ctx sdk.Context, address sdk.Address) int64 {
	store := ctx.KVStore(mapper.key)
	b := store.Get(append([]byte("answer_num-"), address...))
	if len(b) == 0 || len(b) != 8 {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(b))
}

func (mapper BaseQuestionMapper) IncreaseAnswerNum(ctx sdk.Context, address sdk.Address) {
	store := ctx.KVStore(mapper.key)
	i := mapper.GetAnswerNum(ctx, address)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i+1))
	store.Set(append([]byte("answer_num-"), address...), b)
}

func (mapper BaseQuestionMapper) encodeQuestion(question Question) []byte {
	bz, err := mapper.cdc.MarshalBinary(question)
	if err != nil {
		panic(err)
	}
	return bz
}

func (mapper BaseQuestionMapper) decodeQuestion(bz []byte) Question {
	r, n, err := bytes.NewBuffer(bz), new(int), new(error)
	questionI := oldwire.ReadBinary(struct{ Question }{}, r, len(bz), n, err)
	if *err != nil {
		panic(*err)
	}

	question := questionI.(struct{ Question }).Question
	return question
}

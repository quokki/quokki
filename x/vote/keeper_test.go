package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoting(t *testing.T) {
	ctx, vk, _, am, qm, _ := createTestInput(t, false)
	result := qm.CreateQuestion(ctx, addrs[0], addrs[0], "", "", "", "", []string{})
	assert.Equal(t, true, result.IsOK())

	question := result.Data
	result = am.CreateAnswer(ctx, question, addrs[1], "")
	assert.Equal(t, true, result.IsOK())

	answer1 := result.Data
	vk.VoteUp(ctx, addrs[2], answer1, 100)
	vk.VoteUp(ctx, addrs[3], answer1, 10)
	vk.VoteUp(ctx, addrs[4], answer1, -1)

	result = am.CreateAnswer(ctx, question, addrs[2], "")
	assert.Equal(t, true, result.IsOK())

	answer2 := result.Data
	vk.VoteUp(ctx, addrs[2], answer2, 10)
	vk.VoteUp(ctx, addrs[3], answer2, 5)

	assert.Equal(t, qm.GetAnswerNum(ctx, question), int64(2))

	assert.Equal(t, vk.GetQuestionTotalVote(ctx, question), int64(125))
	assert.Equal(t, vk.GetAnswerVote(ctx, answer1), int64(110))
	assert.Equal(t, vk.GetAnswerVote(ctx, answer2), int64(15))
}

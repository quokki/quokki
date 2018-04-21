package commands

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/quokki/quokki/util"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func voteShowCmd(voteStoreName, questionStoreName string, cdc *wire.Codec) *cobra.Command {
	cmdr := showCommander{voteStoreName, questionStoreName, cdc}
	cmd := &cobra.Command{
		Use:   "show <question address>",
		Short: "vote show question and answer",
		RunE:  cmdr.voteShowRun,
	}
	return cmd
}

type showCommander struct {
	voteStoreName     string
	questionStoreName string
	cdc               *wire.Codec
}

func (c showCommander) voteShowRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args[0]) < 1 {
		return errors.New("Need question address")
	}

	addr := args[0]
	bz, err := hex.DecodeString(addr)
	if err != nil {
		return err
	}
	key := sdk.Address(bz)

	res, err := builder.Query(append([]byte("question-total-vote"), key...), c.voteStoreName)
	if err != nil {
		return err
	}
	var questionTotalVote int64 = 0
	err = c.cdc.UnmarshalBinary(res, &questionTotalVote)
	if err != nil {
		return err
	}
	fmt.Printf("Question total vote: %d\n", questionTotalVote)

	res, err = builder.Query(append([]byte("answer_num-"), key...), c.questionStoreName)
	if err != nil {
		return err
	}

	numAnswer := int64(binary.LittleEndian.Uint64(res))

	var i int64
	for i = 0; i < numAnswer; i++ {
		var answerVote int64 = 0
		addr := util.GetAddressIndexHash(key, i, "answer")
		res, err = builder.Query(append([]byte("answer-vote"), addr...), c.voteStoreName)
		if err == nil {
			err := c.cdc.UnmarshalBinary(res, &answerVote)
			if err != nil {
				answerVote = 0
			}
		}
		fmt.Printf("Answer <%s> vote: %d\n", addr.String(), answerVote)
	}

	return nil
}

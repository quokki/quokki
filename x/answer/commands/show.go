package commands

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/quokki/quokki/x/answer"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/quokki/quokki/util"
)

func showAnswerCmd(answerStoreName string, questionStoreName string, cdc *wire.Codec) *cobra.Command {
	cmdr := showCommander{
		answerStoreName,
		questionStoreName,
		cdc,
		answer.GetAnswerDecoder(cdc),
	}
	return &cobra.Command{
		Use:   "show <question address>",
		Short: "Show answer",
		RunE:  cmdr.showAnswerCmd,
	}
}

type showCommander struct {
	answerStoreName   string
	questionStoreName string
	cdc               *wire.Codec
	parser            answer.AnswerDecoder
}

func (c showCommander) showAnswerCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide an address")
	}

	// find the key to look up the account
	addr := args[0]
	bz, err := hex.DecodeString(addr)
	if err != nil {
		return err
	}
	key := sdk.Address(bz)

	res, err := builder.Query(append([]byte("answer_num-"), key...), c.questionStoreName)
	if err != nil {
		return err
	}

	numAnswer := int64(binary.LittleEndian.Uint64(res))

	var i int64
	for i = 0; i < numAnswer; i++ {
		res, err = builder.Query(util.GetAddressIndexHash(key, i, "answer"), c.answerStoreName)
		if err != nil {
			return err
		}
		ans, err := c.parser(res)
		if err != nil {
			return err
		}
		fmt.Println("Writer:", strings.ToUpper(hex.EncodeToString(ans.GetWriter())), "CreateBlock:", ans.GetCreateBlockHeight())

		res, err = builder.Query(append([]byte("content-"), util.GetAddressIndexHash(key, i, "answer")...), c.answerStoreName)
		if err != nil {
			return err
		}
		fmt.Println(string(res))
	}

	return nil
}

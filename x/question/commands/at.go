package commands

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func atQuestionCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := atCommander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "at <block height>",
		Short: "At question",
		RunE:  cmdr.atQuestionCmd,
	}
}

type atCommander struct {
	storeName string
	cdc       *wire.Codec
}

func (c atCommander) atQuestionCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide an block height")
	}

	blockHeight, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(blockHeight))
	res, err := builder.Query(append([]byte("question-at"), b...), c.storeName)
	if err != nil {
		return err
	}

	questions := []sdk.Address{}
	if res != nil && len(res) > 0 {
		err := c.cdc.UnmarshalBinary(res, &questions)
		if err != nil {
			questions = []sdk.Address{}
		}
	}

	for i := 0; i < len(questions); i++ {
		fmt.Println(questions[i].String())
	}
	return nil
}

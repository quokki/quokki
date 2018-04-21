package commands

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/quokki/quokki/x/comment"
)

func createCommentCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := createCommander{cdc}
	cmd := &cobra.Command{
		Use:   "create <target address> <type> <content>",
		Short: "Create comment",
		RunE:  cmdr.createCommentRun,
	}
	return cmd
}

type createCommander struct {
	cdc *wire.Codec
}

func (c createCommander) createCommentRun(cmd *cobra.Command, args []string) error {
	if len(args) < 3 || len(args[0]) < 1 || len(args[1]) < 1 || len(args[2]) < 1 {
		return errors.New("Need target address and type and content")
	}

	// get the from address
	from, err := builder.GetFromAddress()
	if err != nil {
		return err
	}

	// get account name
	name := viper.GetString(client.FlagName)

	// build message
	msg, err := BuildMsg(from, args[0], args[1], args[2])
	if err != nil {
		return err
	}

	// build and sign the transaction, then broadcast to Tendermint
	res, err := builder.SignBuildBroadcast(name, msg, c.cdc)
	if err != nil {
		return err
	}

	fmt.Println("Comment address: ", strings.ToUpper(hex.EncodeToString(res.DeliverTx.Data)))
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func BuildMsg(from sdk.Address, target string, _type string, content string) (sdk.Msg, error) {
	targetAddress, err := hex.DecodeString(target)
	if err != nil {
		return nil, err
	}

	msg := comment.NewCreateCommentMsg(_type, targetAddress, from, content)
	return msg, nil
}

package commands

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/quokki/quokki/x/vote"
)

func voteUpCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := upCommander{cdc}
	cmd := &cobra.Command{
		Use:   "up <answer address> <quokki power>",
		Short: "vote up answer",
		RunE:  cmdr.voteUpRun,
	}
	return cmd
}

type upCommander struct {
	cdc *wire.Codec
}

func (c upCommander) voteUpRun(cmd *cobra.Command, args []string) error {
	if len(args) < 2 || len(args[0]) < 1 || len(args[1]) < 1 {
		return errors.New("Need answer address and quokki power")
	}

	// get the from address
	from, err := builder.GetFromAddress()
	if err != nil {
		return err
	}

	// get account name
	name := viper.GetString(client.FlagName)

	// build message
	msg, err := c.BuildMsg(from, args[0], args[1])
	if err != nil {
		return err
	}

	// build and sign the transaction, then broadcast to Tendermint
	res, err := builder.SignBuildBroadcast(name, msg, c.cdc)
	if err != nil {
		return err
	}

	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func (c upCommander) BuildMsg(from sdk.Address, sAnswer string, sQuokkiPower string) (sdk.Msg, error) {
	answer, err := hex.DecodeString(sAnswer)
	if err != nil {
		return nil, err
	}
	quokkiPower, err := strconv.ParseInt(sQuokkiPower, 10, 64)
	if err != nil {
		return nil, err
	}

	msg := vote.NewVoteUpMsg(from, answer, quokkiPower)
	return msg, nil
}

package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/quokki/quokki/x/power"
)

func powerDownCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := downCommander{cdc}
	cmd := &cobra.Command{
		Use:   "down <power>",
		Short: "Power down T.T",
		RunE:  cmdr.powerDownRun,
	}
	return cmd
}

type downCommander struct {
	cdc *wire.Codec
}

func (c downCommander) powerDownRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args[0]) < 1 {
		return errors.New("Need power amount")
	}

	// get the from address
	from, err := builder.GetFromAddress()
	if err != nil {
		return err
	}

	// get account name
	name := viper.GetString(client.FlagName)

	// build message
	msg, err := c.BuildMsg(from, args[0])
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

func (c downCommander) BuildMsg(from sdk.Address, sPower string) (sdk.Msg, error) {
	pow, err := strconv.ParseInt(sPower, 10, 64)
	if err != nil {
		return nil, err
	}

	msg := power.NewPowerDownMsg(from, pow)
	return msg, nil
}

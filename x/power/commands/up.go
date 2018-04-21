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

func powerUpCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := upCommander{cdc}
	cmd := &cobra.Command{
		Use:   "up <quokki>",
		Short: "Power up Quokki!!!!!",
		RunE:  cmdr.powerUpRun,
	}
	return cmd
}

type upCommander struct {
	cdc *wire.Codec
}

func (c upCommander) powerUpRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args[0]) < 1 {
		return errors.New("Need quokki amount")
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

func (c upCommander) BuildMsg(from sdk.Address, sQuokki string) (sdk.Msg, error) {
	quokki, err := strconv.ParseInt(sQuokki, 10, 64)
	if err != nil {
		return nil, err
	}

	msg := power.NewPowerUpMsg(from, quokki)
	return msg, nil
}

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

func powerUseCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := useCommander{cdc}
	cmd := &cobra.Command{
		Use:   "use <quokki> <restore term>",
		Short: "Test purpose",
		RunE:  cmdr.powerUseRun,
	}
	return cmd
}

type useCommander struct {
	cdc *wire.Codec
}

func (c useCommander) powerUseRun(cmd *cobra.Command, args []string) error {
	if len(args) < 2 || len(args[0]) < 1 || len(args[1]) < 1 {
		return errors.New("Need quokki amount and restore term")
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

func (c useCommander) BuildMsg(from sdk.Address, sQuokki string, sTerm string) (sdk.Msg, error) {
	quokki, err := strconv.ParseInt(sQuokki, 10, 64)
	if err != nil {
		return nil, err
	}
	restoreTerm, err := strconv.ParseInt(sTerm, 10, 64)
	if err != nil {
		return nil, err
	}

	msg := power.NewPowerUseMsg(from, quokki, restoreTerm)
	return msg, nil
}

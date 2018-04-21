package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/quokki/quokki/x/profile"
)

// SendTxCommand will create a send tx and sign it with the given key
func setProfileCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := setCommander{cdc}
	cmd := &cobra.Command{
		Use:   "set <json data>",
		Short: "Create and sign a send tx",
		RunE:  cmdr.setProfileRun,
	}
	return cmd
}

type setCommander struct {
	cdc *wire.Codec
}

func (c setCommander) setProfileRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args[0]) < 1 {
		return errors.New("Need json data")
	}

	// get the from address
	from, err := builder.GetFromAddress()
	if err != nil {
		return err
	}

	// get account name
	name := viper.GetString(client.FlagName)

	// build message
	msg, err := BuildMsg(from, args[0])
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

func BuildMsg(from sdk.Address, data string) (sdk.Msg, error) {
	uprf := profile.UserProfile{}
	err := json.Unmarshal([]byte(data), &uprf)
	if err != nil {
		return nil, err
	}
	msg := profile.NewProfileMsg(from, uprf)
	return msg, nil
}

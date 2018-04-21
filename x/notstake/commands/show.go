package commands

import (
	"fmt"

	"github.com/quokki/quokki/x/notstake"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	"github.com/cosmos/cosmos-sdk/wire"
)

func showNotstakeCmd(notstakeStoreName string, cdc *wire.Codec) *cobra.Command {
	cmdr := showCommander{
		notstakeStoreName,
		cdc,
	}
	return &cobra.Command{
		Use:   "show",
		Short: "Show notstake",
		RunE:  cmdr.showNotstakeCmd,
	}
}

type showCommander struct {
	notstakeStoreName string
	cdc               *wire.Codec
}

func (c showCommander) showNotstakeCmd(cmd *cobra.Command, args []string) error {
	res, err := builder.Query([]byte("notstake-val-infos"), c.notstakeStoreName)
	if err != nil {
		return err
	}

	infos := []notstake.ValInfo{}
	if len(res) > 0 {
		err := c.cdc.UnmarshalBinary(res, &infos)
		if err != nil {
			return err
		}
	}

	for _, info := range infos {
		fmt.Println(info.PubKey.KeyString(), info.Power)
	}

	return nil
}

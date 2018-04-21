package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands to interact with
// local private key storage.
func Commands(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notstake",
		Short: "related notstake",
		Long:  "related notstake",
	}
	setCmd := setNotstakeCmd(cdc)
	client.PostCommands(setCmd)
	showCmd := showNotstakeCmd("main", cdc)
	client.GetCommands(showCmd)
	cmd.AddCommand(
		setCmd,
		showCmd,
	)
	return cmd
}

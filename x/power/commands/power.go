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
		Use:   "power",
		Short: "related power",
		Long:  "related power",
	}
	upCmd := powerUpCmd(cdc)
	client.PostCommands(upCmd)
	downCmd := powerDownCmd(cdc)
	client.PostCommands(downCmd)
	useCmd := powerUseCmd(cdc)
	client.PostCommands(useCmd)
	cmd.AddCommand(
		upCmd,
		downCmd,
		useCmd,
	)
	return cmd
}

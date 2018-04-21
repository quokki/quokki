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
		Use:   "vote",
		Short: "related power",
		Long:  "related power",
	}
	upCmd := voteUpCmd(cdc)
	client.PostCommands(upCmd)
	showCmd := voteShowCmd("vote", "question", cdc)
	client.PostCommands(showCmd)
	cmd.AddCommand(
		upCmd,
		showCmd,
	)
	return cmd
}

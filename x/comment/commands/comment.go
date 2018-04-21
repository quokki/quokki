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
		Use:   "comment",
		Short: "related comment",
		Long:  "related comment",
	}
	createCmd := createCommentCmd(cdc)
	client.PostCommands(createCmd)
	showCmd := showCommentCmd("comment", cdc)
	client.GetCommands(showCmd)
	cmd.AddCommand(
		createCmd,
		showCmd,
	)
	return cmd
}

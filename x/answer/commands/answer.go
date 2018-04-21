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
		Use:   "answer",
		Short: "related answer",
		Long:  "related answer",
	}
	createCmd := createAnswerCmd(cdc)
	client.PostCommands(createCmd)
	updateCmd := updateAnswerCmd(cdc)
	client.PostCommands(updateCmd)
	showCmd := showAnswerCmd("answer", "question", cdc)
	client.GetCommands(showCmd)
	cmd.AddCommand(
		createCmd,
		updateCmd,
		showCmd,
	)
	return cmd
}

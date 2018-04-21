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
		Use:   "question",
		Short: "related question",
		Long:  "related question",
	}
	createCmd := createQuestionCmd(cdc)
	client.PostCommands(createCmd)
	updateCmd := updateQuestionCmd(cdc)
	client.PostCommands(updateCmd)
	showCmd := showQuestionCmd("question", cdc)
	client.GetCommands(showCmd)
	atCmd := atQuestionCmd("question", cdc)
	client.GetCommands(atCmd)
	cmd.AddCommand(
		createCmd,
		updateCmd,
		showCmd,
		atCmd,
	)
	return cmd
}

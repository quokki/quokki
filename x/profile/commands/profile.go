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
		Use:   "profile",
		Short: "Set or get account's profile",
		Long:  `Profile shows your identity.`,
	}
	setCmd := setProfileCmd(cdc)
	client.PostCommands(setCmd)
	getCmd := getProfileCmd("profile", cdc)
	client.GetCommands(getCmd)
	cmd.AddCommand(
		setCmd,
		getCmd,
	)
	return cmd
}

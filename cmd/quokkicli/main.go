package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tmlibs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/commands"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/commands"
	notstakecmd "github.com/quokki/quokki/x/notstake/commands"
	powercmd "github.com/quokki/quokki/x/power/commands"

	answercmd "github.com/quokki/quokki/x/answer/commands"
	commentcmd "github.com/quokki/quokki/x/comment/commands"
	profilecmd "github.com/quokki/quokki/x/profile/commands"
	questioncmd "github.com/quokki/quokki/x/question/commands"
	votecmd "github.com/quokki/quokki/x/vote/commands"

	"github.com/quokki/quokki/app"
	"github.com/quokki/quokki/types"
)

// gaiacliCmd is the entry point for this binary
var (
	basecliCmd = &cobra.Command{
		Use:   "quokkicli",
		Short: "Quokki light-client",
	}
)

func todoNotImplemented(_ *cobra.Command, _ []string) error {
	return errors.New("TODO: Command not yet implemented")
}

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// add standard rpc, and tx commands
	rpc.AddCommands(basecliCmd)
	basecliCmd.AddCommand(client.LineBreak)
	tx.AddCommands(basecliCmd, cdc)
	basecliCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	basecliCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, types.GetAccountDecoder(cdc)),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)
	/*basecliCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCRelayCmd(cdc),
		)...)*/
	// add proxy, version and key info
	basecliCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
		notstakecmd.Commands(cdc),
		powercmd.Commands(cdc),
		profilecmd.Commands(cdc),
		questioncmd.Commands(cdc),
		answercmd.Commands(cdc),
		commentcmd.Commands(cdc),
		votecmd.Commands(cdc),
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(basecliCmd, "BC", os.ExpandEnv("$HOME/.quokkicli"))
	executor.Execute()
}

package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/quokki/quokki/app"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "quokkid",
		Short:             "Quokki Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	viper.SetDefault(server.FlagChainID, "test-chain-0")
	server.AddCommands(ctx, cdc, rootCmd, app.QuokkiAppInit(),
		server.ConstructAppCreator(newApp, "quokki"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "quokki"))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewQuokkiApp(logger, db, traceStore, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	qApp := app.NewQuokkiApp(logger, db, traceStore)
	return qApp.ExportAppStateAndValidators()
}

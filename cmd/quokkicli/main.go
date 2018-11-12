package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

	articlecmd "github.com/quokki/quokki/x/article/client/cli"

	articlerest "github.com/quokki/quokki/x/article/client/rest"

	"github.com/quokki/quokki/app"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "quokkicli",
		Short: "Quokki light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// add standard rpc commands
	rpc.AddCommands(rootCmd)

	//Add state commands
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint state querying subcommands",
	}
	tendermintCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(tendermintCmd, cdc)

	viper.SetDefault(client.FlagChainID, "test-chain-0")
	//Add auth and bank commands
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, authcmd.GetAccountDecoder(cdc)),
			GetQueryCmd(cdc, authcmd.GetAccountDecoder(cdc)),
		)...)
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
			articlecmd.WriteTxCmd(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
func GetQueryCmd(cdc *wire.Codec, decoder auth.AccountDecoder) *cobra.Command {
	return &cobra.Command{
		Use:   "rest-server",
		Short: "Turn on rest api",
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime.GOMAXPROCS(runtime.NumCPU())

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(decoder)

			r := mux.NewRouter()
			authrest.RegisterRoutes(cliCtx, r, cdc, "acc")
			articlerest.RegisterRoutes(cliCtx, r, cdc, "article")
			//http.Handle("/", r)
			go func() {
				if err := http.ListenAndServe(":8080", handlers.CORS()(r)); err != nil {
					log.Println(err)
				}
			}()

			c := make(chan os.Signal, 1)
			// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
			// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
			signal.Notify(c, os.Interrupt)

			// Block until we receive our signal.
			<-c

			log.Println("shutting down")
			os.Exit(0)

			return nil
		},
	}
}

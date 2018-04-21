package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/quokki/quokki/app"
)

// basecoindCmd is the entry point for this binary
var (
	basecoindCmd = &cobra.Command{
		Use:   "quokkid",
		Short: "Quokki Daemon (server)",
	}
)

// defaultOptions sets up the app_options for the
// default genesis file
func defaultOptions(args []string) (json.RawMessage, string, cmn.HexBytes, error) {
	addr, secret, err := server.GenerateCoinKey()
	if err != nil {
		return nil, "", nil, err
	}
	opts := fmt.Sprintf(`{
      "accounts": [{
        "address": "%s",
        "coins": [
          {
            "denom": "quokki",
            "amount": 10000000000000000
          }
        ]
      }],
			"admins": [
				"%s"
			],
			"notstake_tick_param": {
	      "total_notstake_supply": {
	        "num": 2000000000000000,
	        "denom": 1
	      },
				"unused_supply": {
					"num": 0,
					"denom": 1
				},
	      "inflation_rate": {
	        "num": 15,
	        "denom": 100
				}
      },
			"vote_tick_param": {
	      "total_vote_supply": {
	        "num": 7000000000000000,
	        "denom": 1
	      },
				"unused_supply": {
					"num": 0,
					"denom": 1
				},
	      "inflation_rate": {
	        "num": 15,
	        "denom": 100
				}
      },
			"power_tick_param": {
				"total_power_supply": {
	        "num": 1000000000000000,
	        "denom": 1
	      },
				"unused_supply": {
					"num": 0,
					"denom": 1
				},
	      "inflation_rate": {
	        "num": 15,
	        "denom": 100
				}
			}
    }`, addr, addr)
	return json.RawMessage(opts), secret, addr, nil

}

func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	dbMain, err := dbm.NewGoLevelDB("quokki", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbAcc, err := dbm.NewGoLevelDB("quokki-acc", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbIBC, err := dbm.NewGoLevelDB("quokki-ibc", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbPower, err := dbm.NewGoLevelDB("quokki-power", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbVote, err := dbm.NewGoLevelDB("quokki-vote", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbProfile, err := dbm.NewGoLevelDB("quokki-profile", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbQuestion, err := dbm.NewGoLevelDB("quokki-question", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbAnswer, err := dbm.NewGoLevelDB("quokki-answer", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbComment, err := dbm.NewGoLevelDB("quokki-comment", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbs := map[string]dbm.DB{
		"main":     dbMain,
		"acc":      dbAcc,
		"ibc":      dbIBC,
		"power":    dbPower,
		"vote":     dbVote,
		"profile":  dbProfile,
		"question": dbQuestion,
		"answer":   dbAnswer,
		"comment":  dbComment,
	}
	bapp := app.NewBasecoinApp(logger, dbs)
	return bapp, nil
}

func main() {
	// TODO: set logger through CLI
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
		With("module", "main")

	startCmd := server.StartCmd(generateApp, logger)
	startCmd.Flags().String("mgoURL", "", "")
	startCmd.Flags().Bool("mgoTLS", false, "")
	basecoindCmd.AddCommand(
		server.InitCmd(defaultOptions, logger),
		startCmd,
		server.UnsafeResetAllCmd(logger),
		server.ShowNodeIdCmd(logger),
		server.ShowValidatorCmd(logger),
		version.VersionCmd,
	)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.quokkid")
	executor := cli.PrepareBaseCmd(basecoindCmd, "BC", rootDir)
	executor.Execute()
}

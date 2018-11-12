package app

import (
	"encoding/json"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	codec "github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/quokki/quokki/x/article"
	"github.com/quokki/quokki/x/faucet"
)

const (
	appName = "QuokkiApp"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.quokkicli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.quokkid")
)

// Extended ABCI application
type QuokkiApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyArticle       *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	articleKeeper       article.Keeper
	faucetKeeper        faucet.Keeper
}

// NewQuokkiApp returns a reference to an initialized GaiaApp.
func NewQuokkiApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *QuokkiApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	var app = &QuokkiApp{
		BaseApp:          bApp,
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyFeeCollection: sdk.NewKVStoreKey("fee"),
		keyArticle:       sdk.NewKVStoreKey("article"),
	}

	// define the accountMapper
	app.accountMapper = auth.NewAccountMapper(
		app.cdc,
		app.keyAccount,        // target store
		auth.ProtoBaseAccount, // prototype
	)

	// add handlers
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.cdc, app.keyFeeCollection)
	app.articleKeeper = article.NewKeeper(app.cdc, app.keyArticle, app.RegisterCodespace(article.DefaultCodespace))
	app.faucetKeeper = faucet.NewKeeper(app.cdc, app.accountMapper)

	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("article", article.NewHandler(app.articleKeeper)).
		AddRoute("faucet", faucet.NewHandler(app.faucetKeeper))

	// initialize BaseApp
	app.SetInitChainer(app.initChainer)
	anteHandler := auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper)
	app.SetAnteHandler(func(
		ctx sdk.Context, tx sdk.Tx,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {
		msgs := tx.GetMsgs()
		if len(msgs) > 0 {
			msg := msgs[0]
			if msg.Type() == "faucet" {
				newCtx = ctx
				res = sdk.Result{}
				abort = false
				return
			}
		}
		return anteHandler(ctx, tx)
	})
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyFeeCollection, app.keyArticle)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.NewCodec()
	bank.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	article.RegisterCodec(cdc)
	faucet.RegisterCodec(cdc)
	sdk.RegisterWire(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// custom logic for gaia initialization
func (app *QuokkiApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	// TODO is this now the whole genesis file?

	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}

	app.articleKeeper.InitGenesis(ctx, genesisState.GenesisArticle)

	return abci.ResponseInitChain{}
}

// export the state of gaia for a genesis file
func (app *QuokkiApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})

	// iterate to get the accounts
	accounts := []GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)

	genState := GenesisState{
		Accounts: accounts,
	}
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	return appState, nil, nil
}

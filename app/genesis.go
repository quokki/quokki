package app

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	codec "github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/spf13/pflag"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/quokki/quokki/x/article"
)

// DefaultKeyPass contains the default key password for genesis transactions
const DefaultKeyPass = "12345678"

// State to Unmarshal
type GenesisState struct {
	Accounts       []GenesisAccount       `json:"accounts"`
	GenesisArticle article.GenesisArticle `json:"genesisArticle"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
	}
}

func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address: acc.GetAddress(),
		Coins:   acc.GetCoins(),
	}
}

// convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
}

// get app init parameters for server init command
func QuokkiAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(server.FlagName, "", "validator moniker, required")
	fsAppGenTx.String(server.FlagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(server.FlagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         QuokkiAppGenTx,
		AppGenState:      QuokkiAppGenStateJSON,
	}
}

// simple genesis tx
type QuokkiGenTx struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
	PubKey  string         `json:"pub_key"`
}

func QuokkiAppGenTx(
	cdc *codec.Codec, pk crypto.PubKey, genTxConfig config.GenTx,
) (appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	if genTxConfig.Name == "" {
		return nil, nil, tmtypes.GenesisValidator{}, errors.New("Must specify --name (validator moniker)")
	}

	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Password for account '%s' (default %s):", genTxConfig.Name, DefaultKeyPass)

	keyPass, err := client.GetPassword(prompt, buf)
	if err != nil && keyPass != "" {
		// An error was returned that either failed to read the password from
		// STDIN or the given password is not empty but failed to meet minimum
		// length requirements.
		return appGenTx, cliPrint, validator, err
	}

	if keyPass == "" {
		keyPass = DefaultKeyPass
	}

	addr, secret, err := server.GenerateSaveCoinKey(
		genTxConfig.CliRoot,
		genTxConfig.Name,
		keyPass,
		genTxConfig.Overwrite,
	)
	if err != nil {
		return appGenTx, cliPrint, validator, err
	}

	mm := map[string]string{"secret": secret}
	bz, err := cdc.MarshalJSON(mm)
	if err != nil {
		return appGenTx, cliPrint, validator, err
	}

	cliPrint = json.RawMessage(bz)
	appGenTx, _, validator, err = QuokkiAppGenTxNF(cdc, pk, addr, genTxConfig.Name)

	return appGenTx, cliPrint, validator, err
}

// Generate a gaia genesis transaction without flags
func QuokkiAppGenTxNF(cdc *codec.Codec, pk crypto.PubKey, addr sdk.AccAddress, name string) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	var bz []byte
	quokkiGenTx := QuokkiGenTx{
		Name:    name,
		Address: addr,
		PubKey:  sdk.MustBech32ifyAccPub(pk),
	}
	bz, err = codec.MarshalJSONIndent(cdc, quokkiGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  10,
	}
	return
}

// Create the core parameters for genesis initialization for quokki
// note that the pubkey input is this machines pubkey
func QuokkiAppGenState(cdc *codec.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {

	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	genArticle := article.GenesisArticle{}
	// get genesis flag account information
	genaccs := make([]GenesisAccount, len(appGenTxs))
	for i, appGenTx := range appGenTxs {

		var genTx QuokkiGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}

		// create the genesis account, give'm few steaks and a buncha token with there name
		accAuth := auth.NewBaseAccountWithAddress(genTx.Address)
		accAuth.Coins = sdk.Coins{
			{"quokki", sdk.NewInt(10000000000)},
		}
		acc := NewGenesisAccount(&accAuth)
		genaccs[i] = acc

		if i == 0 {
			genArticle.Writer = genTx.Address
			genArticle.Payload = "This is genesis article"
		}
	}

	// create the final app state
	genesisState = GenesisState{
		Accounts:       genaccs,
		GenesisArticle: genArticle,
	}
	return
}

// QuokkiAppGenState but with JSON
func QuokkiAppGenStateJSON(cdc *codec.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	// create the final app state
	genesisState, err := QuokkiAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = codec.MarshalJSONIndent(cdc, genesisState)
	return
}

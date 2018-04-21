package notstake

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/abci/types"

	crypto "github.com/tendermint/go-crypto"
	oldwire "github.com/tendermint/go-wire"

	dbm "github.com/tendermint/tmlibs/db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/quokki/quokki/types"
)

// dummy addresses used for testing
var (
	addrs = []sdk.Address{
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6160"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6161"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6162"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6163"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6164"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6165"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6166"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6167"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6168"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6169"),
	}

	// dummy pubkeys used for testing
	pks = []crypto.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB53"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB54"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB55"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB56"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB57"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB58"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB59"),
	}

	emptyAddr   sdk.Address
	emptyPubkey crypto.PubKey
)

func defaultParams() notstakeTickParam {
	return notstakeTickParam{
		TotalNotstakeSupply: sdk.NewRat(10000000000000000, 1),
		InflationRate:       sdk.NewRat(1, 10),
	}
}

func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, NotstakeKeeper, bank.CoinKeeper) {
	db := dbm.NewMemDB()
	keyMain := sdk.NewKVStoreKey("test")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyMain, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test"}, isCheckTx, nil)
	makeCodec()
	accountMapper := auth.NewAccountMapper(keyMain, &types.AppAccount{})
	coinKeeper := bank.NewCoinKeeper(accountMapper)
	notstakeKeeper := NewNotstakeKeeper(keyMain, coinKeeper)
	notstakeKeeper.SetNotstakeTickParam(ctx, defaultParams())

	for i, pk := range pks {
		notstakeKeeper.SetValInfo(ctx, pk, 1, int64(i+1))
	}

	return ctx, notstakeKeeper, coinKeeper
}

func makeCodec() *wire.Codec {
	const msgTypeNotstake = 0xF
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Msg }{},
		oldwire.ConcreteType{SetMsg{}, msgTypeNotstake},
	)

	const accTypeApp = 0x1
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Account }{},
		oldwire.ConcreteType{&types.AppAccount{}, accTypeApp},
	)
	cdc := wire.NewCodec()

	// cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	// bank.RegisterWire(cdc)   // Register bank.[SendMsg,IssueMsg] types.
	// crypto.RegisterWire(cdc) // Register crypto.[PubKey,PrivKey,Signature] types.
	// ibc.RegisterWire(cdc) // Register ibc.[IBCTransferMsg, IBCReceiveMsg] types.
	return cdc
}

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd crypto.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd.Wrap()
}

func testAddr(addr string) sdk.Address {
	res, err := sdk.GetAddress(addr)
	if err != nil {
		panic(err)
	}
	return res
}

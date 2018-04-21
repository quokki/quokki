package vote

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
	"github.com/quokki/quokki/x/answer"
	"github.com/quokki/quokki/x/power"
	"github.com/quokki/quokki/x/question"
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

func defaultParams() voteTickParam {
	return voteTickParam{
		TotalVoteSupply: sdk.NewRat(10000000000000000, 1),
		UnusedSupply:    sdk.NewRat(0),
		InflationRate:   sdk.NewRat(1, 10),
	}
}

func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, VoteKeeper, power.PowerKeeper, answer.AnswerMapper, question.QuestionMapper, bank.CoinKeeper) {
	db := dbm.NewMemDB()
	keyMain := sdk.NewKVStoreKey("test")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyMain, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test"}, isCheckTx, nil)
	cdc := makeCodec()
	accountMapper := auth.NewAccountMapper(keyMain, &types.AppAccount{})
	coinKeeper := bank.NewCoinKeeper(accountMapper)
	powerKeeper := power.NewPowerKeeper(keyMain, cdc, accountMapper, coinKeeper)
	questionMapper := question.NewQuestionMapper(keyMain)
	answerMapper := answer.NewAnswerMapper(keyMain, questionMapper)
	voteKeeper := NewVoteKeeper(keyMain, cdc, powerKeeper, answerMapper, questionMapper, coinKeeper)
	voteKeeper.SetVoteTickParam(ctx, defaultParams())

	for _, addr := range addrs {
		coinKeeper.AddCoins(ctx, addr, sdk.Coins{
			sdk.Coin{Denom: "quokki", Amount: defaultParams().TotalVoteSupply.Evaluate() / int64(len(addrs))},
		})
		powerKeeper.PowerUp(ctx, addr, (defaultParams().TotalVoteSupply.Evaluate()/int64(len(addrs)))/2)
	}

	return ctx, voteKeeper, powerKeeper, answerMapper, questionMapper, coinKeeper
}

func makeCodec() *wire.Codec {
	const msgTypeSend = 0x1
	const msgTypeIssue = 0x2
	const msgTypeCreateQuestion = 0x8
	const msgTypeCreateAnswer = 0x9
	const msgTypePowerUp = 0xB
	const msgTypePowerDown = 0xC
	const msgTypePowerUse = 0xD
	const msgTypeVoteUp = 0xE
	const msgTypeUpdateQuestion = 0x10
	const msgTypeUpdateAnswer = 0x11
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Msg }{},
		oldwire.ConcreteType{bank.SendMsg{}, msgTypeSend},
		oldwire.ConcreteType{bank.IssueMsg{}, msgTypeIssue},
		oldwire.ConcreteType{question.CreateQuestionMsg{}, msgTypeCreateQuestion},
		oldwire.ConcreteType{answer.CreateAnswerMsg{}, msgTypeCreateAnswer},
		oldwire.ConcreteType{power.PowerUpMsg{}, msgTypePowerUp},
		oldwire.ConcreteType{power.PowerDownMsg{}, msgTypePowerDown},
		oldwire.ConcreteType{power.PowerUseMsg{}, msgTypePowerUse},
		oldwire.ConcreteType{VoteUpMsg{}, msgTypeVoteUp},
		oldwire.ConcreteType{question.UpdateQuestionMsg{}, msgTypeUpdateQuestion},
		oldwire.ConcreteType{answer.UpdateAnswerMsg{}, msgTypeUpdateAnswer},
	)

	const accTypeApp = 0x1
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Account }{},
		oldwire.ConcreteType{&types.AppAccount{}, accTypeApp},
	)
	cdc := wire.NewCodec()

	question.RegisterWire()
	answer.RegisterWire()
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

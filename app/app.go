package app

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"

	"github.com/spf13/viper"

	abci "github.com/tendermint/abci/types"
	oldwire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"

	"github.com/quokki/quokki/db"

	"github.com/quokki/quokki/types"
	"github.com/quokki/quokki/util"
	"github.com/quokki/quokki/x/answer"
	"github.com/quokki/quokki/x/comment"
	"github.com/quokki/quokki/x/notstake"
	"github.com/quokki/quokki/x/power"
	"github.com/quokki/quokki/x/profile"
	"github.com/quokki/quokki/x/question"
	register "github.com/quokki/quokki/x/simpleregister"
	"github.com/quokki/quokki/x/vote"

	"github.com/globalsign/mgo"
)

const (
	appName = "BasecoinApp"
)

// Extended ABCI application
type BasecoinApp struct {
	*bam.BaseApp
	cdc         *wire.Codec
	queryRouter util.QueryRouter
	valUpdates  []abci.Validator

	// keys to access the substores
	capKeyMainStore     *sdk.KVStoreKey
	capKeyAccountStore  *sdk.KVStoreKey
	capKeyIBCStore      *sdk.KVStoreKey
	capKeyPowerStore    *sdk.KVStoreKey
	capKeyVoteStore     *sdk.KVStoreKey
	capKeyProfileStore  *sdk.KVStoreKey
	capKeyQuestionStore *sdk.KVStoreKey
	capKeyAnswerStore   *sdk.KVStoreKey
	capKeyCommentStore  *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper  sdk.AccountMapper
	notstakeKeeper notstake.NotstakeKeeper
	powerKeeper    power.PowerKeeper
	voteKeeper     vote.VoteKeeper
	profileMapper  profile.ProfileMapper
	questionMapper question.QuestionMapper
	answerMapper   answer.AnswerMapper
	commentMapper  comment.CommentMapper
}

func NewBasecoinApp(logger log.Logger, dbs map[string]dbm.DB) *BasecoinApp {
	// create your application object
	var app = &BasecoinApp{
		BaseApp:             bam.NewBaseApp(appName, logger, dbs["main"]),
		cdc:                 MakeCodec(),
		queryRouter:         util.NewQueryRouter(),
		capKeyMainStore:     sdk.NewKVStoreKey("main"),
		capKeyAccountStore:  sdk.NewKVStoreKey("acc"),
		capKeyIBCStore:      sdk.NewKVStoreKey("ibc"),
		capKeyPowerStore:    sdk.NewKVStoreKey("power"),
		capKeyVoteStore:     sdk.NewKVStoreKey("vote"),
		capKeyProfileStore:  sdk.NewKVStoreKey("profile"),
		capKeyQuestionStore: sdk.NewKVStoreKey("question"),
		capKeyAnswerStore:   sdk.NewKVStoreKey("answer"),
		capKeyCommentStore:  sdk.NewKVStoreKey("comment"),
	}

	if len(viper.GetString("mgoURL")) > 0 {
		dialInfo, err := mgo.ParseURL(viper.GetString("mgoURL"))
		if viper.GetBool("mgoTLS") {
			tlsConfig := &tls.Config{}
			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
				return conn, err
			}
		}
		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		db.SetSession(session)
	}

	// define the accountMapper
	app.accountMapper = auth.NewAccountMapperSealed(
		app.capKeyAccountStore, // target store
		&types.AppAccount{},    // prototype
	)

	coinKeeper := bank.NewCoinKeeper(app.accountMapper)
	ibcMapper := ibc.NewIBCMapper(app.cdc, app.capKeyIBCStore)

	app.powerKeeper = power.NewPowerKeeper(app.capKeyPowerStore, app.cdc, app.accountMapper, coinKeeper)
	app.notstakeKeeper = notstake.NewNotstakeKeeper(app.capKeyMainStore, coinKeeper)
	app.profileMapper = profile.NewProfileMapper(app.capKeyProfileStore)
	app.questionMapper = question.NewQuestionMapper(app.capKeyQuestionStore)
	app.answerMapper = answer.NewAnswerMapper(app.capKeyAnswerStore, app.questionMapper)
	commentTypeToInfo := comment.CommentTypeToInfo{
		"answer":   comment.CommentInfo{Key: app.capKeyAnswerStore, CollectionName: "answers"},
		"question": comment.CommentInfo{Key: app.capKeyQuestionStore, CollectionName: "questions"},
	}
	app.commentMapper = comment.NewCommentMapper(app.capKeyCommentStore, commentTypeToInfo)

	app.voteKeeper = vote.NewVoteKeeper(app.capKeyVoteStore, app.cdc, app.powerKeeper, app.answerMapper, app.questionMapper, coinKeeper)

	// add handlers
	notstakeHandler := notstake.NewHandler(app.notstakeKeeper)
	app.Router().
		AddRoute("bank", bank.NewHandler(coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(ibcMapper, coinKeeper)).
		AddRoute("notstake", func(ctx sdk.Context, msg sdk.Msg) (result sdk.Result) {
			result = notstakeHandler(ctx, msg)
			if ctx.IsCheckTx() == false && result.IsOK() {
				app.valUpdates = append(app.valUpdates, result.ValidatorUpdates...)
			}
			return
		}).
		AddRoute("power", power.NewHandler(app.powerKeeper)).
		AddRoute("vote", vote.NewHandler(app.voteKeeper)).
		AddRoute("profile", profile.NewHandler(app.profileMapper)).
		AddRoute("question", question.NewHandler(app.questionMapper)).
		AddRoute("answer", answer.NewHandler(app.answerMapper)).
		AddRoute("comment", comment.NewHandler(app.commentMapper)).
		AddRoute("register", register.NewHandler(app.accountMapper))

	app.SetEndBlocker(app.endBlocker)

	app.queryRouter.
		AddRoute("sequence", util.NewSequenceQueryHandler(*app.capKeyAccountStore, types.GetAccountDecoder(app.cdc))).
		AddRoute("account", util.NewAccountQueryHandler(*app.capKeyAccountStore, types.GetAccountDecoder(app.cdc)))

	// initialize BaseApp
	app.SetTxDecoder(app.txDecoder)
	app.SetInitChainer(app.initChainer)
	app.MountStoreWithDB(app.capKeyMainStore, sdk.StoreTypeIAVL, dbs["main"])
	app.MountStoreWithDB(app.capKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	app.MountStoreWithDB(app.capKeyIBCStore, sdk.StoreTypeIAVL, dbs["ibc"])
	app.MountStoreWithDB(app.capKeyPowerStore, sdk.StoreTypeIAVL, dbs["power"])
	app.MountStoreWithDB(app.capKeyVoteStore, sdk.StoreTypeIAVL, dbs["vote"])
	app.MountStoreWithDB(app.capKeyProfileStore, sdk.StoreTypeIAVL, dbs["profile"])
	app.MountStoreWithDB(app.capKeyQuestionStore, sdk.StoreTypeIAVL, dbs["question"])
	app.MountStoreWithDB(app.capKeyAnswerStore, sdk.StoreTypeIAVL, dbs["answer"])
	app.MountStoreWithDB(app.capKeyCommentStore, sdk.StoreTypeIAVL, dbs["comment"])

	anteHandler := auth.NewAnteHandler(app.accountMapper)
	app.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx) (_ sdk.Context, _ sdk.Result, abort bool) {
		msg := tx.GetMsg()
		if msg != nil {
			if msg.Type() == "register" {
				return ctx, sdk.Result{}, false
			}
		}
		return anteHandler(ctx, tx)
	})
	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

func (app *BasecoinApp) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) (response abci.ResponseEndBlock) {
	app.powerKeeper.Tick(ctx)
	app.voteKeeper.Tick(ctx)
	response.ValidatorUpdates = app.valUpdates
	app.valUpdates = nil
	return
}

// custom tx codec
// TODO: use new go-wire
func MakeCodec() *wire.Codec {
	const msgTypeSend = 0x1
	const msgTypeIssue = 0x2
	const msgTypeIBCTransferMsg = 0x5
	const msgTypeIBCReceiveMsg = 0x6
	const msgTypeProfile = 0x7
	const msgTypeCreateQuestion = 0x8
	const msgTypeCreateAnswer = 0x9
	const msgTypeCreateComment = 0xA
	const msgTypePowerUp = 0xB
	const msgTypePowerDown = 0xC
	const msgTypePowerUse = 0xD
	const msgTypeVoteUp = 0xE
	const msgTypeNotstake = 0xF
	const msgTypeUpdateQuestion = 0x10
	const msgTypeUpdateAnswer = 0x11
	const msgTypeRegister = 0x12
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Msg }{},
		oldwire.ConcreteType{bank.SendMsg{}, msgTypeSend},
		oldwire.ConcreteType{bank.IssueMsg{}, msgTypeIssue},
		oldwire.ConcreteType{ibc.IBCTransferMsg{}, msgTypeIBCTransferMsg},
		oldwire.ConcreteType{ibc.IBCReceiveMsg{}, msgTypeIBCReceiveMsg},
		oldwire.ConcreteType{profile.ProfileMsg{}, msgTypeProfile},
		oldwire.ConcreteType{question.CreateQuestionMsg{}, msgTypeCreateQuestion},
		oldwire.ConcreteType{answer.CreateAnswerMsg{}, msgTypeCreateAnswer},
		oldwire.ConcreteType{comment.CreateCommentMsg{}, msgTypeCreateComment},
		oldwire.ConcreteType{power.PowerUpMsg{}, msgTypePowerUp},
		oldwire.ConcreteType{power.PowerDownMsg{}, msgTypePowerDown},
		oldwire.ConcreteType{power.PowerUseMsg{}, msgTypePowerUse},
		oldwire.ConcreteType{vote.VoteUpMsg{}, msgTypeVoteUp},
		oldwire.ConcreteType{notstake.SetMsg{}, msgTypeNotstake},
		oldwire.ConcreteType{question.UpdateQuestionMsg{}, msgTypeUpdateQuestion},
		oldwire.ConcreteType{answer.UpdateAnswerMsg{}, msgTypeUpdateAnswer},
		oldwire.ConcreteType{register.RegisterMsg{}, msgTypeRegister},
	)

	const accTypeApp = 0x1
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Account }{},
		oldwire.ConcreteType{&types.AppAccount{}, accTypeApp},
	)
	cdc := wire.NewCodec()

	profile.RegisterWire()
	question.RegisterWire()
	answer.RegisterWire()
	// cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	// bank.RegisterWire(cdc)   // Register bank.[SendMsg,IssueMsg] types.
	// crypto.RegisterWire(cdc) // Register crypto.[PubKey,PrivKey,Signature] types.
	// ibc.RegisterWire(cdc) // Register ibc.[IBCTransferMsg, IBCReceiveMsg] types.
	return cdc
}

// custom logic for transaction decoding
func (app *BasecoinApp) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = sdk.StdTx{}

	if len(txBytes) == 0 {
		return nil, sdk.ErrTxDecode("txBytes are empty")
	}

	// StdTx.Msg is an interface. The concrete types
	// are registered by MakeTxCodec in bank.RegisterWire.
	var err error
	if txBytes[0] == '{' && txBytes[len(txBytes)-1] == '}' {
		err = app.cdc.UnmarshalJSON(txBytes, &tx)
	} else {
		err = app.cdc.UnmarshalBinary(txBytes, &tx)
	}
	if err != nil {
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}

	return tx, nil
}

// custom logic for basecoin initialization
func (app *BasecoinApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := json.Unmarshal(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	var quokkiSum int64 = 0

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		app.accountMapper.SetAccount(ctx, acc)
		quokkiSum += acc.GetCoins().AmountOf("quokki")
	}

	voteGenesis := new(vote.VoteTickParamGenesis)
	err = json.Unmarshal(stateJSON, voteGenesis)
	if err != nil {
		panic(err)
	}
	app.voteKeeper.SetVoteTickParam(ctx, voteGenesis.Param)
	quokkiSum -= voteGenesis.Param.TotalVoteSupply.Evaluate()

	powerGenesis := new(power.PowerTickParamGenesis)
	err = json.Unmarshal(stateJSON, powerGenesis)
	if err != nil {
		panic(err)
	}
	app.powerKeeper.SetPowerTickParam(ctx, powerGenesis.Param)
	quokkiSum -= powerGenesis.Param.TotalPowerSupply.Evaluate()

	notstakeGenesis := new(notstake.NotstakeTickParamGenesis)
	err = json.Unmarshal(stateJSON, notstakeGenesis)
	if err != nil {
		panic(err)
	}
	app.notstakeKeeper.SetNotstakeTickParam(ctx, notstakeGenesis.Param)
	quokkiSum -= notstakeGenesis.Param.TotalNotstakeSupply.Evaluate()

	if quokkiSum != 0 {
		panic(fmt.Sprintf("Total supply does not match! (gap: %d)", quokkiSum))
	}

	notstakeAdmins := new(notstake.NotstakeAdminGenesis)
	err = json.Unmarshal(stateJSON, notstakeAdmins)
	if err != nil {
		panic(err)
	}
	notstake.SetAdmins(notstakeAdmins.Admins)

	return abci.ResponseInitChain{}
}

func (app *BasecoinApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	defer func() {
		if r := recover(); r != nil {
			log := fmt.Sprintf("Recovered when query", r)
			res = sdk.ErrInternal(log).Result().ToQuery()
		}
	}()

	optionName, subpath, err := util.ParsePath(req.Path)
	if err != nil {
		return err.Result().ToQuery()
	}

	if optionName == "api" {
		optionName, subpath, err := util.ParsePath(subpath)
		if err != nil {
			return err.Result().ToQuery()
		}
		h := app.queryRouter.Route(optionName)
		if h != nil {
			req.Path = subpath
			res = h(app.BaseApp, req)
		} else {
			res = sdk.ErrUnknownRequest("Invalid query: " + optionName).Result().ToQuery()
		}
	} else {
		res = app.BaseApp.Query(req)
	}
	return
}

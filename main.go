package main

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	codec "github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// Cosmos-sdk version = v0.24.2
// amino version = v0.12.0-rc0
func main() {
	cdc := codec.NewCodec()

	codec.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	bank.RegisterWire(cdc)

	issuer, _ := sdk.AccAddressFromBech32("cosmosaccaddr12pga2ll5472pkjjtytntk7pa0chpqhxprex49p") //arbitrary address
	fmt.Println(hex.EncodeToString(issuer))
	issueMsg := bank.NewMsgIssue(
		issuer,
		[]bank.Output{bank.NewOutput(issuer, sdk.Coins{sdk.NewCoin("test", sdk.NewInt(10000))})},
	)

	// StdTx is a standard way to wrap a Msg with Fee and Signatures.
	// NOTE: the first signature is the FeePayer (Signatures must not be nil).
	/*type StdTx struct {
		Msgs       []sdk.Msg      `json:"msg"`	// <- sdk.Msg is interface
		Fee        StdFee         `json:"fee"`
		Signatures []StdSignature `json:"signatures"`
		Memo       string         `json:"memo"`
	}*/
	bz, _ := hex.DecodeString("03e9db07a9a70eafacbbf2d9ecef2e5275074fc39f137cbe7942e84a18ca05cd43")
	var pubKey [33]byte
	copy(pubKey[:], bz[:33])
	stdTx := auth.StdTx{
		Msgs: []sdk.Msg{issueMsg},
		Fee:  auth.NewStdFee(200000, sdk.NewCoin("test", sdk.NewInt(1))),
		Signatures: []auth.StdSignature{auth.StdSignature{
			PubKey:        secp256k1.PubKeySecp256k1(pubKey),
			Signature:     []byte{0},
			AccountNumber: 1,
			Sequence:      1,
		}},
		Memo: "test",
	}

	bz, _ = cdc.MarshalBinary(stdTx)
	fmt.Println(hex.EncodeToString(bz))

	json, _ := cdc.MarshalJSON(stdTx)
	fmt.Println(string(json))

	hex, err := hex.DecodeString("8f01f0625dee0a41c06abad60a145051d57ff4af941b4a4b22e6bb783d7e2e105cc112250a145051d57ff4af941b4a4b22e6bb783d7e2e105cc1120d0a047465737412053130303030120f0a090a04746573741201311080b5181a2f0a26eb5ae9872103e9db07a9a70eafacbbf2d9ecef2e5275074fc39f137cbe7942e84a18ca05cd4312010018022002220474657374")
	if err != nil {
		panic(err)
	}
	stdTx = auth.StdTx{}
	cdc.UnmarshalBinary(hex, &stdTx)
	fmt.Println(stdTx)
}

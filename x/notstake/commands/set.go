package commands

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/quokki/quokki/x/notstake"
	crypto "github.com/tendermint/go-crypto"
)

func setNotstakeCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := setCommander{cdc}
	cmd := &cobra.Command{
		Use:   "set <pub key> <power> <weight>",
		Short: "set",
		RunE:  cmdr.setNotstakeRun,
	}
	return cmd
}

type setCommander struct {
	cdc *wire.Codec
}

func (c setCommander) setNotstakeRun(cmd *cobra.Command, args []string) error {
	if len(args) < 3 || len(args[0]) < 1 || len(args[1]) < 1 || len(args[2]) < 1 {
		return errors.New("Need pub key and power and weight")
	}

	// get the from address
	from, err := builder.GetFromAddress()
	if err != nil {
		return err
	}

	// get account name
	name := viper.GetString(client.FlagName)

	// build message
	msg, err := BuildMsg(from, args[0], args[1], args[2])
	if err != nil {
		return err
	}

	// build and sign the transaction, then broadcast to Tendermint
	res, err := builder.SignBuildBroadcast(name, msg, c.cdc)
	if err != nil {
		return err
	}

	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func BuildMsg(from sdk.Address, sPubkey string, sPower string, sWeight string) (sdk.Msg, error) {
	bPubKey, err := hex.DecodeString(sPubkey)
	if err != nil {
		return nil, err
	}
	if len(bPubKey) != 32 {
		return nil, errors.New("Invalid pub key length")
	}
	pubKey := crypto.PubKeyEd25519{}
	copy(pubKey[:], bPubKey[:32])

	power, err := strconv.ParseInt(sPower, 10, 64)
	if err != nil {
		return nil, err
	}
	weight, err := strconv.ParseInt(sWeight, 10, 64)
	if err != nil {
		return nil, err
	}
	msg := notstake.NewSetMsg(pubKey.Wrap(), power, weight, from)
	return msg, nil
}

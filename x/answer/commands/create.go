package commands

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/quokki/quokki/x/answer"
)

func createAnswerCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := createCommander{cdc}
	cmd := &cobra.Command{
		Use:   "create <question address> <json data>",
		Short: "Create question",
		RunE:  cmdr.createAnswerRun,
	}
	return cmd
}

type createCommander struct {
	cdc *wire.Codec
}

func (c createCommander) createAnswerRun(cmd *cobra.Command, args []string) error {
	if len(args) < 2 || len(args[0]) < 1 || len(args[1]) < 1 {
		return errors.New("Need question address and json data")
	}

	// get the from address
	from, err := builder.GetFromAddress()
	if err != nil {
		return err
	}

	// get account name
	name := viper.GetString(client.FlagName)

	// build message
	msg, err := c.BuildMsg(from, args[0], args[1])
	if err != nil {
		return err
	}

	// build and sign the transaction, then broadcast to Tendermint
	res, err := builder.SignBuildBroadcast(name, msg, c.cdc)
	if err != nil {
		return err
	}

	fmt.Println("Answer address: ", strings.ToUpper(hex.EncodeToString(res.DeliverTx.Data)))
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func (c createCommander) BuildMsg(from sdk.Address, question string, jsonData string) (sdk.Msg, error) {
	questionAddress, err := hex.DecodeString(question)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}
	mapData := data.(map[string]interface{})
	content, ok := mapData["content"].(string)
	if ok == false {
		return nil, errors.New("Invalid content")
	}
	msg := answer.NewCreateAnswerMsg(questionAddress, from, content)
	return msg, nil
}

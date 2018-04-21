package commands

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/quokki/quokki/x/answer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func updateAnswerCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := updateCommander{cdc}
	cmd := &cobra.Command{
		Use:   "update <question address> <json data>",
		Short: "Update question",
		RunE:  cmdr.updateAnswerRun,
	}
	return cmd
}

type updateCommander struct {
	cdc *wire.Codec
}

func (c updateCommander) updateAnswerRun(cmd *cobra.Command, args []string) error {
	if len(args) < 2 || len(args[0]) < 1 || len(args[1]) < 1 {
		return errors.New("Need answer address and json data")
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

func (c updateCommander) BuildMsg(from sdk.Address, sAnswer string, jsonData string) (sdk.Msg, error) {
	answerAddress, err := hex.DecodeString(sAnswer)
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
	msg := answer.NewUpdateAnswerMsg(answerAddress, from, content)
	return msg, nil
}

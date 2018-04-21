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
	"github.com/quokki/quokki/x/question"
)

func createQuestionCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := createCommander{cdc}
	cmd := &cobra.Command{
		Use:   "create <json data>",
		Short: "Create question",
		RunE:  cmdr.createQuestionRun,
	}
	return cmd
}

type createCommander struct {
	cdc *wire.Codec
}

func (c createCommander) createQuestionRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args[0]) < 1 {
		return errors.New("Need json data")
	}

	// get the from address
	from, err := builder.GetFromAddress()
	if err != nil {
		return err
	}

	// get account name
	name := viper.GetString(client.FlagName)

	// build message
	msg, err := c.BuildMsg(from, args[0])
	if err != nil {
		return err
	}

	// build and sign the transaction, then broadcast to Tendermint
	res, err := builder.SignBuildBroadcast(name, msg, c.cdc)
	if err != nil {
		return err
	}

	fmt.Println("Question address: ", strings.ToUpper(hex.EncodeToString(res.DeliverTx.Data)))
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func (c createCommander) BuildMsg(from sdk.Address, jsonData string) (sdk.Msg, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}
	mapData := data.(map[string]interface{})
	title, ok := mapData["title"].(string)
	if ok == false {
		return nil, errors.New("Invalid title")
	}
	content, ok := mapData["content"].(string)
	if ok == false {
		return nil, errors.New("Invalid content")
	}
	language, ok := mapData["language"].(string)
	if ok == false {
		return nil, errors.New("Invalid language")
	}
	category, ok := mapData["category"].(string)
	if ok == false {
		return nil, errors.New("Invalid category")
	}
	_tags, ok := mapData["tags"].([]interface{})
	if ok == false {
		return nil, errors.New("Invalid tags")
	}
	tags := make([]string, 0, len(_tags))
	for _, _tag := range _tags {
		tag, ok := _tag.(string)
		if ok == false {
			return nil, errors.New("Invalid tag")
		}
		tags = append(tags, tag)
	}
	msg := question.NewCreateQuestionMsg(from, title, content, language, category, tags)
	return msg, nil
}

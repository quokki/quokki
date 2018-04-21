package commands

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/quokki/quokki/x/question"
)

func updateQuestionCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := updateCommander{cdc}
	cmd := &cobra.Command{
		Use:   "update <json data>",
		Short: "Update question",
		RunE:  cmdr.updateQuestionRun,
	}
	return cmd
}

type updateCommander struct {
	cdc *wire.Codec
}

func (c updateCommander) updateQuestionRun(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func (c updateCommander) BuildMsg(from sdk.Address, sAddress, jsonData string) (sdk.Msg, error) {
	address, err := hex.DecodeString(sAddress)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal([]byte(jsonData), &data)
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
	msg := question.NewUpdateQuestionMsg(address, from, title, content, language, category, tags)
	return msg, nil
}

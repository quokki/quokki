package commands

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/quokki/quokki/x/question"
)

func showQuestionCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := showCommander{
		storeName,
		cdc,
		question.GetQuestionDecoder(cdc),
	}
	return &cobra.Command{
		Use:   "show <address>",
		Short: "Show question",
		RunE:  cmdr.showQuestionCmd,
	}
}

type showCommander struct {
	storeName string
	cdc       *wire.Codec
	parser    question.QuestionDecoder
}

func (c showCommander) showQuestionCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide an address")
	}

	// find the key to look up the account
	addr := args[0]
	bz, err := hex.DecodeString(addr)
	if err != nil {
		return err
	}
	key := sdk.Address(bz)

	res, err := builder.Query(key, c.storeName)
	if err != nil {
		return err
	}

	// parse out the value
	question, err := c.parser(res)
	if err != nil {
		return err
	}

	// print out whole account
	output, err := json.MarshalIndent(question, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	res, err = builder.Query(append([]byte("title-"), key...), c.storeName)
	if err != nil {
		return err
	}

	fmt.Println("Title: " + string(res))

	res, err = builder.Query(append([]byte("language-"), key...), c.storeName)
	if err == nil {
		fmt.Println("Language: " + string(res))
	}

	res, err = builder.Query(append([]byte("category-"), key...), c.storeName)
	if err == nil {
		fmt.Println("Category: " + string(res))
	}

	res, err = builder.Query(append([]byte("tags-"), key...), c.storeName)
	if err == nil {
		tags := []string{}
		err = c.cdc.UnmarshalBinary(res, &tags)
		if err == nil {
			fmt.Print("Tags: ")
			for _, tag := range tags {
				fmt.Print(tag + " ")
			}
			fmt.Println("")
		}
	}

	res, err = builder.Query(append([]byte("content-"), key...), c.storeName)
	if err != nil {
		return err
	}

	fmt.Println("Content: " + string(res))

	return nil
}

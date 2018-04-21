package commands

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/quokki/quokki/util"
	"github.com/quokki/quokki/x/comment"
)

func showCommentCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := showCommander{
		storeName,
		cdc,
		comment.GetCommentDecoder(cdc),
	}
	return &cobra.Command{
		Use:   "show <target address>",
		Short: "Show comment",
		RunE:  cmdr.showCommentCmd,
	}
}

type showCommander struct {
	storeName string
	cdc       *wire.Codec
	parser    comment.CommentDecoder
}

func (c showCommander) showCommentCmd(cmd *cobra.Command, args []string) error {
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

	res, err := builder.Query(append([]byte("comment_num-"), key...), c.storeName)
	if err != nil {
		return err
	}

	numAnswer := int64(binary.LittleEndian.Uint64(res))

	var i int64
	for i = 0; i < numAnswer; i++ {
		res, err = builder.Query(util.GetAddressIndexHash(key, i, "comment"), c.storeName)
		if err != nil {
			return err
		}
		comt, err := c.parser(res)
		if err != nil {
			return err
		}
		fmt.Println("Writer:", hex.EncodeToString(comt.GetWriter()), "CreateBlock:", comt.GetCreateBlockHeight())
		fmt.Println(comt.GetContent())
	}

	return nil
}

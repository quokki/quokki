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

	"github.com/quokki/quokki/x/profile"
)

func getProfileCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := getCommander{
		storeName,
		cdc,
		profile.GetProfileDecoder(cdc),
	}
	return &cobra.Command{
		Use:   "get <address>",
		Short: "Query account profile",
		RunE:  cmdr.getProfileCmd,
	}
}

type getCommander struct {
	storeName string
	cdc       *wire.Codec
	parser    profile.ProfileDecoder
}

func (c getCommander) getProfileCmd(cmd *cobra.Command, args []string) error {
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

	res, err := builder.Query(append([]byte("profile-"), key...), c.storeName)
	if err != nil {
		return err
	}

	// parse out the value
	profile, err := c.parser(res)
	if err != nil {
		return err
	}

	// print out whole account
	output, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	return nil
}

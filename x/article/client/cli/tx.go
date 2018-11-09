package cli

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	codec "github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/quokki/quokki/x/article"
)

const (
	flagParentId = "parent"
	flagPayload  = "payload"
)

func WriteTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Write article",
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			parentId := viper.GetInt64(flagParentId)
			payload := viper.GetString(flagPayload)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			parent := []byte{}
			if parentId >= 0 {
				bz := make([]byte, 8)
				binary.LittleEndian.PutUint64(bz, uint64(parentId))
				parent = bz
			}
			msg := article.NewMsgWrite(from, parent, payload)

			json, _ := cdc.MarshalJSON(msg)
			fmt.Println("!!!")
			fmt.Println(string(json))
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int64(flagParentId, -1, "Parent id")
	cmd.Flags().String(flagPayload, "", "Payload to inject")

	return cmd
}

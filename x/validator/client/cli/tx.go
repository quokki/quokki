package cli

import (
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

	"github.com/quokki/quokki/x/validator"
)

const (
	flagPubKey = "pubkey"
	flagPower  = "power"
)

func ValidatorTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "set validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			strPubKey := viper.GetString(flagPubKey)
			pubKey, err := sdk.GetValPubKeyBech32(strPubKey)
			if err != nil {
				return err
			}

			power := viper.GetInt64(flagPower)

			valAddr := sdk.ValAddress(pubKey.Address())
			fmt.Println(valAddr.String())

			msg := validator.NewMsgValidator(from, valAddr, pubKey, power)
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagPubKey, "", "Pubkey")
	cmd.Flags().Int64(flagPower, 1, "power")

	return cmd
}

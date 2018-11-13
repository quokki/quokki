package faucet

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	codec "github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var alphabet = [5]string{"A", "B", "C", "D", "F"}

type Keeper struct {
	accountMapper auth.AccountMapper
}

func NewKeeper(cdc *codec.Codec, accountMapper auth.AccountMapper) Keeper {
	return Keeper{
		accountMapper: accountMapper,
	}
}

func (keeper Keeper) NewFaucet(ctx sdk.Context, address sdk.AccAddress) (sdk.Tags, sdk.Error) {
	tags := sdk.EmptyTags()

	account := keeper.accountMapper.GetAccount(ctx, address)
	if account != nil {
		return sdk.EmptyTags(), sdk.ErrInternal("Already existed account")
	}

	account = keeper.accountMapper.NewAccountWithAddress(ctx, address)
	if account == nil {
		return sdk.EmptyTags(), sdk.ErrInternal("Fail to create account")
	}

	al := alphabet[ctx.BlockHeight()%5]

	// send coin and random token for test
	account.SetCoins(sdk.Coins{
		sdk.Coin{
			Denom:  "quokki",
			Amount: sdk.NewInt(100000),
		},
		sdk.Coin{
			Denom:  "test" + al,
			Amount: sdk.NewInt(10000),
		},
	})

	keeper.accountMapper.SetAccount(ctx, account)

	tags.AppendTag("faucet", address)
	return tags, nil
}

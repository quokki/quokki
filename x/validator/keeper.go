package validator

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	codec "github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/tendermint/crypto"
)

type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}

func (keeper Keeper) SetValidator(ctx sdk.Context, address sdk.ValAddress, pubKey crypto.PubKey, power int64) (sdk.Tags, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)

	validator := Validator{
		Address:     address,
		PubKey:      pubKey,
		VotingPower: power,
	}
	bz, err := keeper.cdc.MarshalBinaryBare(validator)
	if err != nil {
		return sdk.EmptyTags(), sdk.ErrInternal(err.Error())
	}

	store.Set([]byte("validator"), bz)

	return sdk.EmptyTags(), nil
}

func (keeper Keeper) ClearValidator(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set([]byte("validator"), []byte{})
}

func (keeper Keeper) GetValidator(ctx sdk.Context) (Validator, sdk.Error) {
	store := ctx.TransientStore(keeper.storeKey)

	bz := store.Get([]byte("validator"))
	if len(bz) == 0 {
		return Validator{}, nil
	}

	validator := Validator{}
	err := keeper.cdc.UnmarshalBinaryBare(bz, &validator)
	if err != nil {
		return Validator{}, sdk.ErrInternal(err.Error())
	}

	return validator, nil
}

func (keeper Keeper) SetAdmin(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set([]byte("admin"), address)
}

func (keeper Keeper) IsValidAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(keeper.storeKey)
	admin := store.Get([]byte("admin"))

	if bytes.Compare(admin, address) == 0 {
		return true
	}

	return false
}

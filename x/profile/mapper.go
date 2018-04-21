package profile

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"

	"github.com/quokki/quokki/db"
)

var _ ProfileMapper = (*profileMapper)(nil)

type profileMapper struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewProfileMapper(key sdk.StoreKey) profileMapper {
	cdc := wire.NewCodec()
	return profileMapper{
		key: key,
		cdc: cdc,
	}
}

func (pm profileMapper) GetProfile(ctx sdk.Context, address sdk.Address) Profile {
	store := ctx.KVStore(pm.key)
	bz := store.Get(append([]byte("profile-"), address...))
	if bz == nil {
		return nil
	}
	profile := pm.DecodeAccount(bz)
	return profile
}

func (pm profileMapper) SetProfile(ctx sdk.Context, profile Profile, address sdk.Address) sdk.Result {
	store := ctx.KVStore(pm.key)
	bz := pm.EncodeAccount(profile)
	store.Set(append([]byte("profile-"), address...), bz)

	db.Upsert(ctx, "profiles", map[string]interface{}{"address": address.String()}, profile, map[string]interface{}{"address": address.String()})
	return sdk.Result{}
}

func (pm profileMapper) GetStoreName() string {
	return pm.key.Name()
}

func (pm profileMapper) EncodeAccount(profile Profile) []byte {
	bz, err := pm.cdc.MarshalBinary(profile)
	if err != nil {
		panic(err)
	}
	return bz
}

func (pm profileMapper) DecodeAccount(bz []byte) Profile {
	var profile = UserProfile{}

	err := pm.cdc.UnmarshalBinary(bz, &profile)
	if err != nil {
		return nil
	}
	return &profile
}

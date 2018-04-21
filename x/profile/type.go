package profile

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"
)

type Profile interface {
	ValidateBasic() sdk.Error
	GetSignBytes() []byte
}

type ProfileMapper interface {
	GetProfile(ctx sdk.Context, addr sdk.Address) Profile
	SetProfile(ctx sdk.Context, profile Profile, account sdk.Address) sdk.Result

	EncodeAccount(profile Profile) []byte
	DecodeAccount(bz []byte) Profile
}

func RegisterWire() {
	const prfTypeUser = 0x1
	var _ = oldwire.RegisterInterface(
		struct{ Profile }{},
		oldwire.ConcreteType{&UserProfile{}, prfTypeUser},
	)
}

type ProfileDecoder func([]byte) (Profile, error)

func GetProfileDecoder(cdc *wire.Codec) ProfileDecoder {
	return func(prfBytes []byte) (Profile, error) {
		var profile = UserProfile{}

		err := cdc.UnmarshalBinary(prfBytes, &profile)
		return &profile, err
	}
}

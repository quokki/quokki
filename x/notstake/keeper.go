package notstake

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	crypto "github.com/tendermint/go-crypto"
)

/*
If people are approved by the administrator,
people can become a validator without any risk at present.
TODO: If tendermint has matured enough, I will implement DPOS.
*/

type NotstakeKeeper struct {
	ck bank.CoinKeeper

	key sdk.StoreKey
	cdc *wire.Codec
}

func NewNotstakeKeeper(key sdk.StoreKey, coinKeeper bank.CoinKeeper) NotstakeKeeper {
	cdc := wire.NewCodec()
	return NotstakeKeeper{
		key: key,
		cdc: cdc,
		ck:  coinKeeper,
	}
}

func (keeper NotstakeKeeper) GetValInfos(ctx sdk.Context) []ValInfo {
	infos := []ValInfo{}
	store := ctx.KVStore(keeper.key)
	bz := store.Get([]byte("notstake-val-infos"))
	if bz != nil && len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &infos)
		if err != nil {
			infos = []ValInfo{}
		}
	}
	return infos
}

func (keeper NotstakeKeeper) SetValInfo(ctx sdk.Context, pubKey crypto.PubKey, power int64, weight int64) error {
	infos := keeper.GetValInfos(ctx)
	infoIndex := -1
	for i, info := range infos {
		if bytes.Equal(info.PubKey.Bytes(), pubKey.Bytes()) {
			infoIndex = i
			infos[i].Power = power
			infos[i].Weight = weight
		}
	}
	if power <= 0 && infoIndex >= 0 {
		infos = append(infos[0:infoIndex], infos[infoIndex+1:]...)
	}

	if power > 0 && infoIndex < 0 {
		info := ValInfo{PubKey: pubKey, Power: power, Weight: weight}
		infos = append(infos, info)
	}

	store := ctx.KVStore(keeper.key)
	bz, err := keeper.cdc.MarshalBinary(infos)
	if err != nil {
		return err
	}
	store.Set([]byte("notstake-val-infos"), bz)
	return nil
}

package power

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/quokki/quokki/types"
)

type PowerKeeper struct {
	key sdk.StoreKey
	cdc *wire.Codec
	am  sdk.AccountMapper
	ck  bank.CoinKeeper
}

func NewPowerKeeper(key sdk.StoreKey, cdc *wire.Codec, am sdk.AccountMapper, ck bank.CoinKeeper) PowerKeeper {
	return PowerKeeper{key: key, cdc: cdc, am: am, ck: ck}
}

func (keeper PowerKeeper) PowerUp(ctx sdk.Context, address sdk.Address, quokki int64) sdk.Result {
	account := keeper.am.GetAccount(ctx, address)
	if account == nil {
		return sdk.ErrUnknownAddress(address.String()).Result()
	}

	if quokki <= 0 {
		return sdk.ErrInvalidCoins("Should not be zero").Result()
	}

	coins := account.GetCoins()
	newCoins := coins.Minus(sdk.Coins{sdk.Coin{Amount: quokki, Denom: "quokki"}})
	if !newCoins.IsNotNegative() {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("%d - %d = %d", coins.AmountOf("quokki"), quokki, newCoins.AmountOf("quokki"))).Result()
	}

	err := account.SetCoins(newCoins)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	_quokkiPower, err := account.Get("QuokkiPower")
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	quokkiPower, ok := _quokkiPower.(types.QuokkiPower)
	quokkiPower.Available += quokki
	if ok == false {
		return sdk.ErrInternal("Fail to cast type").Result()
	}

	err = account.Set("QuokkiPower", quokkiPower)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	keeper.am.SetAccount(ctx, account)
	return sdk.Result{}
}

func (keeper PowerKeeper) PowerDown(ctx sdk.Context, address sdk.Address, quokkiPower int64) sdk.Result {
	account := keeper.am.GetAccount(ctx, address)
	if account == nil {
		return sdk.ErrUnknownAddress(address.String()).Result()
	}

	_quokkiPower, err := account.Get("QuokkiPower")
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	accountQuokkiPower, ok := _quokkiPower.(types.QuokkiPower)
	if ok == false {
		return sdk.ErrInternal("Fail to cast type").Result()
	}

	resultQuokkiPower := accountQuokkiPower
	resultQuokkiPower.Available -= quokkiPower
	if resultQuokkiPower.Available < 0 {
		return sdk.ErrInsufficientFunds(fmt.Sprintf("%d - %d = %d", accountQuokkiPower.Available, quokkiPower, resultQuokkiPower.Available)).Result()
	}
	err = account.Set("QuokkiPower", resultQuokkiPower)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	plusCoin := sdk.Coin{Amount: quokkiPower, Denom: "quokki"}
	coins := account.GetCoins()
	newCoins := coins.Plus(sdk.Coins{plusCoin})

	err = account.SetCoins(newCoins)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	keeper.am.SetAccount(ctx, account)

	return sdk.Result{}
}

func (keeper PowerKeeper) PowerUse(ctx sdk.Context, address sdk.Address, quokkiPower int64, restoreTerm int64) error {
	if restoreTerm <= 0 {
		return sdk.ErrInternal("restore term should be positive")
	}

	account := keeper.am.GetAccount(ctx, address)
	if account == nil {
		return sdk.ErrUnknownAddress(address.String())
	}

	_quokkiPower, err := account.Get("QuokkiPower")
	if err != nil {
		return err
	}
	accountQuokkiPower, ok := _quokkiPower.(types.QuokkiPower)
	if ok == false {
		return sdk.ErrInternal("Fail to cast type")
	}

	resultQuokkiPower := accountQuokkiPower
	resultQuokkiPower.Available -= quokkiPower
	resultQuokkiPower.Used += quokkiPower
	if resultQuokkiPower.Available < 0 {
		return sdk.ErrInsufficientFunds(fmt.Sprintf("%d - %d = %d", accountQuokkiPower.Available, quokkiPower, resultQuokkiPower.Available))
	}
	err = account.Set("QuokkiPower", resultQuokkiPower)
	if err != nil {
		return err
	}
	if ctx.IsCheckTx() {
		return nil
	}

	restoreBlockHeight := ctx.BlockHeight() + restoreTerm
	store := ctx.KVStore(keeper.key)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(restoreBlockHeight))
	dest := append([]byte("restore-"), b...)
	bz := store.Get(dest)

	restoreInfos := []RestoreInfo{}
	if len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &restoreInfos)
		if err != nil {
			return err
		}
	}

	restoreInfo := RestoreInfo{Address: address, QuokkiPower: quokkiPower}
	restoreInfos = append(restoreInfos, restoreInfo)
	bz, err = keeper.cdc.MarshalBinary(restoreInfos)
	if err != nil {
		return err
	}
	store.Set(dest, bz)
	keeper.am.SetAccount(ctx, account)

	inflationBlockHeight := ((restoreBlockHeight / BlocksPerProvision) + 1) * BlocksPerProvision
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(inflationBlockHeight))
	dest = append([]byte("inflation-"), b...)
	bz = store.Get(dest)

	restoreInfos = []RestoreInfo{}
	if len(bz) > 0 {
		err := keeper.cdc.UnmarshalBinary(bz, &restoreInfos)
		if err != nil {
			return err
		}
	}

	restoreInfos = append(restoreInfos, restoreInfo)
	bz, err = keeper.cdc.MarshalBinary(restoreInfos)
	if err != nil {
		return err
	}
	store.Set(dest, bz)

	return nil
}

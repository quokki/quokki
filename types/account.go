package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ sdk.Account = (*AppAccount)(nil)

// Custom extensions for this application.  This is just an example of
// extending auth.BaseAccount with custom fields.
//
// This is compatible with the stock auth.AccountStore, since
// auth.AccountStore uses the flexible go-wire library.
type AppAccount struct {
	auth.BaseAccount

	QuokkiPower QuokkiPower `json:"quokki_power"`
}

// Get the ParseAccount function for the custom AppAccount
func GetAccountDecoder(cdc *wire.Codec) sdk.AccountDecoder {
	return func(accBytes []byte) (res sdk.Account, err error) {
		acct := new(AppAccount)
		err = cdc.UnmarshalBinary(accBytes, &acct)
		if err != nil {
			panic(err)
		}
		return acct, err
	}
}

func (account AppAccount) Get(key interface{}) (value interface{}, err error) {
	if key == "QuokkiPower" {
		return account.QuokkiPower, nil
	}
	return nil, nil
}

func (account *AppAccount) Set(key interface{}, value interface{}) error {
	if key == "QuokkiPower" {
		power, ok := value.(QuokkiPower)
		if ok {
			account.QuokkiPower = power
			return nil
		}
	}
	return errors.New("Invalid key or type")
}

//___________________________________________________________________________________

// State to Unmarshal
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.Address `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
}

func NewGenesisAccount(aa *AppAccount) *GenesisAccount {
	return &GenesisAccount{
		Address: aa.Address,
		Coins:   aa.Coins,
	}
}

// convert GenesisAccount to AppAccount
func (ga *GenesisAccount) ToAppAccount() (acc *AppAccount, err error) {
	baseAcc := auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins,
	}
	return &AppAccount{
		BaseAccount: baseAcc,
	}, nil
}

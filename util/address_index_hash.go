package util

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"golang.org/x/crypto/ripemd160"
)

/*
Make new address with address.
If index is same, it will cause conflict address problemn.
Make sure that index is not same.
*/
func GetAddressIndexHash(address sdk.Address, index int64, _type string) sdk.Address {
	iBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(iBytes, uint64(index))
	bytes := make([]byte, 0, 28)
	bytes = append(bytes, address...)
	bytes = append(bytes, iBytes...)
	bytes = append(bytes, _type...)
	h := ripemd160.New()
	h.Write(bytes)
	return h.Sum(nil)
}

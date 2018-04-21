package util

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"

	"github.com/quokki/quokki/types"
)

func NewAccountQueryHandler(keyStore sdk.KVStoreKey, decoder sdk.AccountDecoder) types.QueryHandler {
	return func(baseapp *bam.BaseApp, req abci.RequestQuery) (res abci.ResponseQuery) {
		addr := req.Path
		if len(addr) <= 1 {
			return sdk.ErrUnknownRequest("You should specify address.").Result().ToQuery()
		}
		if addr[0] == '/' {
			addr = string(addr[1:])
		}
		i := strings.Index(addr, "/")
		if i >= 0 {
			if i == len(addr)-1 {
				addr = string(addr[0 : len(addr)-1])
			} else {
				return sdk.ErrUnknownRequest("Invalid path.").Result().ToQuery()
			}
		}

		dst := make([]byte, hex.DecodedLen(len(addr)))
		_, err := hex.Decode(dst, []byte(addr))
		if err != nil {
			return sdk.ErrUnknownRequest(err.Error()).Result().ToQuery()
		}
		req.Data = dst
		req.Path = "/" + keyStore.Name() + "/store"
		req.Prove = false
		res = baseapp.Query(req)
		if res.Value == nil && len(res.Value) == 0 {
			return
		}
		account, err := decoder(res.Value)
		if err != nil {
			return sdk.ErrUnknownRequest(err.Error()).Result().ToQuery()
		}
		res.Value, err = json.Marshal(account)
		if err != nil {
			return sdk.ErrUnknownRequest(err.Error()).Result().ToQuery()
		}
		return
	}
}

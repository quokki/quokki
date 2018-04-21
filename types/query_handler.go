package types

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	abci "github.com/tendermint/abci/types"
)

type QueryHandler func(baseapp *bam.BaseApp, req abci.RequestQuery) (res abci.ResponseQuery)

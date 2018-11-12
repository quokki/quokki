package article

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Article struct {
	Id        []byte         `json:"id"`
	Writer    sdk.AccAddress `json:"writer"`
	Parent    []byte         `json:"parent"`
	Sequence  uint64         `json:"sequence"` // Similar with number of children
	CreatedAt time.Time      `json:"createdAt"`
	Payload   string         `json:"payload"` // TODO: Only save hash of payload
}

// Genesis article has empty id
// So genesis article will be parent of articles that have no explicit parent
type GenesisArticle struct {
	Writer  sdk.AccAddress `json:"writer"`
	Payload string         `json:"payload"`
}

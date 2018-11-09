package article

import sdk "github.com/cosmos/cosmos-sdk/types"

type Article struct {
	Writer  sdk.AccAddress `json:"writer"`
	Parent  []byte         `json:"parent"`
	Payload []byte         `json:"payload"` //TODO: Only save hash of payload
}

func NewArticle(writer sdk.AccAddress, parent []byte, payload []byte) Article {
	return Article{
		Writer:  writer,
		Parent:  parent,
		Payload: payload,
	}
}

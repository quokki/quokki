package article

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 10

	CodeTooBigPayload      sdk.CodeType = 101
	CodeNonexistentArticle sdk.CodeType = 102
	CodeInvalidArticle     sdk.CodeType = 103
	CodeAssignedArticle    sdk.CodeType = 104
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeTooBigPayload:
		return "Too big payload"
	case CodeNonexistentArticle:
		return "Non existent article"
	case CodeInvalidArticle:
		return "Invalid article"
	case CodeAssignedArticle:
		return "Already assigned article"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

func ErrTooBigPayload(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTooBigPayload, codeToDefaultMsg(CodeTooBigPayload))
}

func ErrNonexistentArticle(codespace sdk.CodespaceType, id []byte) sdk.Error {
	return sdk.NewError(codespace, CodeNonexistentArticle, fmt.Sprintf("%s: %s", codeToDefaultMsg(CodeNonexistentArticle), hex.EncodeToString(id)))
}

func ErrInvalidArticle(codespace sdk.CodespaceType, id []byte) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArticle, fmt.Sprintf("%s: %s", codeToDefaultMsg(CodeInvalidArticle), hex.EncodeToString(id)))
}

func ErrAssignedArticle(codespace sdk.CodespaceType, id []byte) sdk.Error {
	return sdk.NewError(codespace, CodeAssignedArticle, fmt.Sprintf("%s: %s", codeToDefaultMsg(CodeAssignedArticle), hex.EncodeToString(id)))
}

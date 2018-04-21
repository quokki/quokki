package util

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ParsePath expects a format like /<optioneName>[/<subpath>]
// Must start with /, subpath may be empty
// Returns error if it doesn't start with /
func ParsePath(path string) (optionName string, subpath string, err sdk.Error) {
	if !strings.HasPrefix(path, "/") {
		err = sdk.ErrUnknownRequest(fmt.Sprintf("invalid path: %s", path))
		return
	}
	paths := strings.SplitN(path[1:], "/", 2)
	optionName = paths[0]
	if len(paths) == 2 {
		subpath = "/" + paths[1]
	}
	return
}

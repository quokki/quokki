package profile

import (
	"encoding/json"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ Profile = (*UserProfile)(nil)

type UserProfile struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	About    string `json:"about"`
}

func (profile UserProfile) ValidateBasic() sdk.Error {
	err := ""
	if utf8.RuneCountInString(profile.Name) > 20 || len(profile.Name) > 80 {
		err += "Name must be shorter than 20 characters."
	}
	if utf8.RuneCountInString(profile.Email) > 30 || len(profile.Email) > 120 {
		err += "Email must be shorter than 30 characters."
	}
	if utf8.RuneCountInString(profile.Phone) > 20 || len(profile.Phone) > 80 {
		err += "Phone must be shorter than 20 characters."
	}
	if utf8.RuneCountInString(profile.Location) > 100 || len(profile.Location) > 400 {
		err += "Location must be shorter than 100 characters."
	}
	if utf8.RuneCountInString(profile.About) > 400 || len(profile.About) > 1600 {
		err += "About must be shorter than 400 characters."
	}
	if err != "" {
		return sdk.NewError(301, err)
	}
	return nil
}

func (profile UserProfile) GetSignBytes() []byte {
	b, err := json.Marshal(profile) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

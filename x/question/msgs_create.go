package question

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MIN_TITLE_LENGTH = 5
const MAX_TITLE_LENGTH = 200
const MAX_CONTENT_LENGTH = 50000
const MAX_LANGUAGE_LENGTH = 30
const MAX_CATEGORY_LENGTH = 50
const MAX_TAG_LENGTH = 20
const MAX_TAGS_LENGTH = 10

type CreateQuestionMsg struct {
	Writer   sdk.Address `json:"writer"`
	Title    string      `json:"title"`
	Content  string      `json:"content"`
	Language string      `json:"language"`
	Category string      `json:"category"`
	Tags     []string    `json:"tags"`
}

var _ sdk.Msg = CreateQuestionMsg{}

func validateBasicQuestionMsg(title string, content string, language string, category string, tags []string) sdk.Error {
	err := ""
	if utf8.RuneCountInString(title) < MIN_TITLE_LENGTH {
		err += fmt.Sprintf("Title must be longer than %d characters.", MIN_TITLE_LENGTH)
	}
	if utf8.RuneCountInString(title) > MAX_TITLE_LENGTH || len(title) > MAX_TITLE_LENGTH*4 {
		err += fmt.Sprintf("Title must be shorter than %d characters.", MAX_TITLE_LENGTH)
	}
	if utf8.RuneCountInString(content) > MAX_CONTENT_LENGTH || len(content) > MAX_CONTENT_LENGTH*4 {
		err += fmt.Sprintf("Content must be shorter than %d characters.", MAX_CONTENT_LENGTH)
	}
	if len(language) == 0 {
		err += fmt.Sprint("You should set language. See IETF language tag and BCP47")
	}
	if utf8.RuneCountInString(language) > MAX_LANGUAGE_LENGTH || len(language) > MAX_LANGUAGE_LENGTH*4 {
		err += fmt.Sprintf("Langauge must be shorter than %d characters.", MAX_LANGUAGE_LENGTH)
	}
	if utf8.RuneCountInString(category) > MAX_CATEGORY_LENGTH || len(category) > MAX_CATEGORY_LENGTH*4 {
		err += fmt.Sprintf("Category must be shorter than %d characters.", MAX_CATEGORY_LENGTH)
	}
	if len(tags) > MAX_TAGS_LENGTH {
		err += fmt.Sprintf("Tags should be fewer than %d.", MAX_TAGS_LENGTH)
	}
	for _, tag := range tags {
		if utf8.RuneCountInString(tag) > MAX_TAG_LENGTH || len(tag) > MAX_TAG_LENGTH*4 {
			err += fmt.Sprintf("Tag must be shorter than %d characters.", MAX_TAG_LENGTH)
		}
	}
	if err != "" {
		return sdk.NewError(401, err)
	}
	return nil
}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewCreateQuestionMsg(writer sdk.Address, title string, content string, language string, category string, tags []string) CreateQuestionMsg {
	return CreateQuestionMsg{Writer: writer, Title: title, Content: content, Language: language, Category: category, Tags: tags}
}

// Implements Msg.
func (msg CreateQuestionMsg) Type() string { return "question" }

// Implements Msg.
func (msg CreateQuestionMsg) ValidateBasic() sdk.Error {
	if len(msg.Writer) != 20 {
		return sdk.ErrInvalidAddress("Invalid address")
	}

	return validateBasicQuestionMsg(msg.Title, msg.Content, msg.Language, msg.Category, msg.Tags)
}

func (msg CreateQuestionMsg) String() string {
	return fmt.Sprintf("CreateQuestionMsg{%v}", msg.Writer)
}

// Implements Msg.
func (msg CreateQuestionMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg CreateQuestionMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg CreateQuestionMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Writer}
}

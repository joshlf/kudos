package purple_unicorn

import (
	"fmt"
	"regexp"
)

type Validator interface {
	Validate() error
	MustValidate()
}

type Code string

var re = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9_]+")

func (c Code) Validate() error {
	if len(c) == 0 {
		return fmt.Errorf("code must be nonempty")
	}
	if !re.MatchString(string(c)) {
		return fmt.Errorf("string contains illegal characters: must be alphanumeric and start with an alphabetic character")
	}
	return nil
}

func (c Code) MustValidate() {
	err := c.Validate()
	if err != nil {
		panic(err)
	}
}

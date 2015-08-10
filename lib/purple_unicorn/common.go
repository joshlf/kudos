package purple_unicorn

import (
	"fmt"
	"regexp"
	"time"
)

// Validator is a type which can
// validate that its state conforms
// to whatever constraints are
// associated with that type.
type Validator interface {
	// Validate performs validation, returning
	// an error if the internal state is invalid.
	// Implementations of the Validator interface
	// provided by this package will often return
	// a value which is of type ErrList.
	Validate() error

	// MustValidate calls Validate, and panics
	// with its return value if it is not nil.
	MustValidate()
}

type Code string

var re = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9_]*")

// Validate implements the Validator Validate method.
func (c Code) Validate() error {
	if len(c) == 0 {
		return fmt.Errorf("must be nonempty")
	}
	if !re.MatchString(string(c)) {
		return fmt.Errorf("contains illegal characters; must be alphanumeric and start with an alphabetic character")
	}
	return nil
}

// MustValidate implements the Validator MustValidate method.
func (c Code) MustValidate() {
	err := c.Validate()
	if err != nil {
		panic(err)
	}
}

// pretty-format an error message about c
func codeErrMsg(c Code, err error) error {
	if c == "" {
		return fmt.Errorf("bad code: %v", err)
	}
	return fmt.Errorf("bad code %q: %v", c, err)
}

type date time.Time

func (d *date) UnmarshalText(text []byte) error {
	t, err := timeparse(string(text))
	if err != nil {
		return err
	}
	*d = date(t)
	return nil
}

func timeparse(text string) (time.Time, error) {
	return time.Parse("Jan 2, 2006 at 3:04pm (MST)", text)
}

func mustPanic(err error, fn string) {
	panic(fmt.Errorf("purple_unicorn: %v: %v", fn, err))
}

package purple_unicorn

import (
	"fmt"
	"regexp"
	"strings"
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

var re = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9_]+")

// Validate implements the Validator Validate method.
func (c Code) Validate() error {
	if len(c) == 0 {
		return fmt.Errorf("code must be nonempty")
	}
	if !re.MatchString(string(c)) {
		return fmt.Errorf("string contains illegal characters: must be alphanumeric and start with an alphabetic character")
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

func mustPanic(err error, fn string, args ...interface{}) {
	percentVs := make([]string, len(args))
	for i := range percentVs {
		percentVs[i] = "%v"
	}
	fmtString := "purple_unicorn: %v(" + strings.Join(percentVs, ", ") + "): %v"
	errorfArgs := []interface{}{fn}
	errorfArgs = append(errorfArgs, args...)
	errorfArgs = append(errorfArgs, err)
	panic(fmt.Errorf(fmtString, errorfArgs...))
}

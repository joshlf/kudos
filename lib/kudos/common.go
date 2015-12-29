package kudos

import (
	"fmt"
	"regexp"
	"time"
)

var re = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9_]*")

func ValidateCode(code string) error {
	switch {
	case len(code) == 0:
		return fmt.Errorf("must be non-empty")
	case !re.MatchString(code):
		return fmt.Errorf("contains illegal characters; must be alphanumeric and start with an alphabetic character")
	}
	return nil
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

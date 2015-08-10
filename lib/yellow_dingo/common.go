package yellow_dingo

import (
	"time"
)

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

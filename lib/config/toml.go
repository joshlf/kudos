package config

import (
	"fmt"
	"strings"
	"time"
)

// number represents a number which is
// either a floating point or a decimal
type number float64

func (n *number) UnmarshalTOML(i interface{}) error {
	f, ok := i.(float64)
	if !ok {
		ii, ok := i.(int64)
		if !ok {
			return fmt.Errorf("expected number")
		}
		f = float64(ii)
	}
	*n = number(f)
	return nil
}

func (n number) MarshalText() ([]byte, error) {
	// TODO(synful): This may be the behavior
	// of fmt.Sprint anyway, so we can get rid
	// of this method entirely
	if number(int(n)) == n {
		return []byte(fmt.Sprint(int(n))), nil
	}
	return []byte(fmt.Sprint(n)), nil
}

// This type should be used whenever
// a safe code is needed (course code,
// assignment code, problem code, etc)
type code string

func (c *code) UnmarshalTOML(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected string value")
	}
	if len(str) == 0 {
		return fmt.Errorf("cannot be empty string")
	}
	const special = string("#+-:@^_")
	const allowedChars = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVQXYZ0123456789" + special
	for _, c := range str {
		if strings.Index(allowedChars, string(c)) == -1 {
			return fmt.Errorf("contains illegal characters; must be alphanumeric or one of %v", special)
		}
	}
	*c = code(str)
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

type optionalDate struct {
	date date
	set  bool
}

func (o *optionalDate) UnmarshalText(text []byte) error {
	o.set = true
	return o.date.UnmarshalText(text)
}

type optionalNumber struct {
	number number
	set    bool
}

func (o *optionalNumber) UnmarshalTOML(i interface{}) error {
	o.set = true
	return o.number.UnmarshalTOML(i)
}

type optionalString struct {
	string string
	set    bool
}

func (o *optionalString) UnmarshalTOML(i interface{}) error {
	o.set = true
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected string")
	}
	o.string = s
	return nil
}

type optionalCode struct {
	code code
	set  bool
}

func (o *optionalCode) UnmarshalTOML(i interface{}) error {
	o.set = true
	return o.code.UnmarshalTOML(i)
}

type optionalHandinMethod struct {
	handinMethod handinMethod
	set          bool
}

func (o *optionalHandinMethod) UnmarshalTOML(i interface{}) error {
	o.set = true
	return o.handinMethod.UnmarshalTOML(i)
}

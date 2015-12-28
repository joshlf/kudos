package config

import (
	"testing"
)

func TestTOML(t *testing.T) {
	var n number
	testError(t, func() error { err := n.UnmarshalTOML("foo"); return err }, "expected number")

	var c code
	var i interface{}
	f := func() error { err := c.UnmarshalTOML(i); return err }
	i = 0
	testError(t, f, "expected string value")
	i = ""
	testError(t, f, "cannot be empty string")
	i = "foo/bar"
	testError(t, f, "contains illegal characters; must be alphanumeric or one of #+-:@^_")

	var o optionalString
	testError(t, func() error { err := o.UnmarshalTOML(0); return err }, "expected string")
}

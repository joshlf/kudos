// Package dev provides functionality which exists
// solely for development purposes, and will eventually
// be removed. In particular, it provides that functionality
// which is used by multiple other packages.
package dev

import (
	"fmt"
	"os"
)

func Fail() {
	fmt.Fprintln(os.Stderr, "[dev] failing for lack of anything better to do")
	os.Exit(1)
}

// Package dev provides functionality which exists
// solely for development purposes, and will eventually
// be removed. In particular, it provides that functionality
// which is used by multiple other packages.
package dev

func Fail() {
	msg := "[dev] failing for lack of anything better to do"
	// if build.DevMode {
	panic(msg)
	// }
	// fmt.Fprintln(os.Stderr, msg)
	// os.Exit(1)
}

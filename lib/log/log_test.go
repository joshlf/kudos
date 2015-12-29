package log

import "fmt"

func ExampleLevels() {
	lgr := NewLogger()
	// In case the test is compiled in debug mode
	lgr.SetLevel(Info)

	// Test that level is Info by default
	lgr.Verbose.Println("initial: not printed")
	lgr.Info.Println("initial: info")

	// Test that each level suppresses
	// output on the levels below it, and
	// use both the lgr.Info.Printf style
	// and the lgr.Printf(level, ...) style.
	funcs := []Printer{lgr.Debug, lgr.Verbose, lgr.Info, lgr.Warn, lgr.Error}
	for l := Debug; l <= Error; l++ {
		lgr.SetLevel(l)
		for ll := Debug; ll <= Error; ll++ {
			funcs[int(ll)].Printf("%v: %v\n", l, ll)
			lgr.Printf(ll, "%v: %v\n", l, ll)
		}
		fmt.Println()
	}

	// Output: initial: info
	// Debug: Debug
	// Debug: Debug
	// Debug: Verbose
	// Debug: Verbose
	// Debug: Info
	// Debug: Info
	// Debug: Warn
	// Debug: Warn
	// Debug: Error
	// Debug: Error
	//
	// Verbose: Verbose
	// Verbose: Verbose
	// Verbose: Info
	// Verbose: Info
	// Verbose: Warn
	// Verbose: Warn
	// Verbose: Error
	// Verbose: Error
	//
	// Info: Info
	// Info: Info
	// Info: Warn
	// Info: Warn
	// Info: Error
	// Info: Error
	//
	// Warn: Warn
	// Warn: Warn
	// Warn: Error
	// Warn: Error
	//
	// Error: Error
	// Error: Error
}

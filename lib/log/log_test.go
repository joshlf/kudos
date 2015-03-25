package log

import "fmt"

func ExampleLevels() {
	oldlevel := level

	// Test that level is Info by default
	Verbose.Println("initial: not printed")
	Info.Println("initial: info")

	// Test that each level suppresses
	// output on the levels below it
	for l := Debug; l <= Error; l++ {
		SetLoggingLevel(l)
		for ll := Debug; ll <= Error; ll++ {
			ll.Printf("%v: %v\n", l, ll)
		}
		fmt.Println()
	}

	level = oldlevel

	// Output: initial: info
	// Debug: Debug
	// Debug: Verbose
	// Debug: Info
	// Debug: Warn
	// Debug: Error
	//
	// Verbose: Verbose
	// Verbose: Info
	// Verbose: Warn
	// Verbose: Error
	//
	// Info: Info
	// Info: Warn
	// Info: Error
	//
	// Warn: Warn
	// Warn: Error
	//
	// Error: Error
}

package perm

import (
	"fmt"
	"os"
)

const (
	Execute os.FileMode = (1 << iota)
	Write
	Read
)

// Parse parses a standard Unix
// permission string (rwxrwxrwx)
// and returns the corresponding
// os.FileMode. perm must be 9
// characters long and be properly
// formatted or else Parse
// will panic.
func Parse(perm string) os.FileMode {
	var mode os.FileMode
	if len(perm) != 9 {
		panic("perm: perm string must be of length 9")
	}
	const on = "rwxrwxrwx"
	const off = "---------"
	for i := 0; i < 9; i++ {
		mode <<= 1
		switch perm[i] {
		case on[i]:
			mode |= 1
		case off[i]:
		default:
			panic(fmt.Errorf("perm: malformed perm string: %v", perm))
		}
	}
	return mode
}

// ParseSingle parses a single component
// of a standard Unix permission string
// (rwx) and returns the corresponding
// os.FileMode. perm must be 3 characters
// long and be properly formatted or else
// Parse will panic.
func ParseSingle(perm string) os.FileMode {
	var mode os.FileMode
	if len(perm) != 3 {
		panic("perm: perm string must be of length 3")
	}
	const on = "rwx"
	const off = "---"
	for i := 0; i < 3; i++ {
		mode <<= 1
		switch perm[i] {
		case on[i]:
			mode |= 1
		case off[i]:
		default:
			panic(fmt.Errorf("perm: malformed perm string: %v", perm))
		}
	}
	return mode
}

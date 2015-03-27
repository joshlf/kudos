package perm

import (
	"fmt"
	"os"
	"os/exec"
)

type Perm uint8

const (
	Execute Perm = (1 << iota)
	Write
	Read
)

func (p Perm) String() string {
	var s string
	c := '-'
	if p&Read != 0 {
		c = 'r'
	}
	s += string(c)
	c = '-'
	if p&Write != 0 {
		c = 'w'
	}
	s += string(c)
	c = '-'
	if p&Execute != 0 {
		c = 'x'
	}
	s += string(c)
	return s
}

type Entity uint8

const (
	User Entity = iota
	Group
	Other
)

func (e Entity) aclCode() string {
	switch e {
	case User:
		return "u"
	case Group:
		return "g"
	case Other:
		return "o"
	default:
		panic(fmt.Errorf("perm: invalid Entity %v", uint8(e)))
	}
}

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

type Facl struct {
	Entity Entity
	Name   string
	Perm   Perm
}

func (f Facl) aclArg() string {
	return f.Entity.aclCode() + f.Name + f.Perm.String()
}

// AddFacl adds the given facls to file.
// It panics if len(facls) < 1.
func AddFacl(file string, facls ...Facl) error {
	// TODO(synful): implement without setfacl
	// dependancy
	// TODO(synful): test
	if len(facls) < 1 {
		panic("perm: no facls provided")
	}
	arg := faclArgString(facls...)
	cmd := exec.Command("setfacl", "-m", arg)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("run setfacl: %v", err)
	}
	return nil
}

// SetFacl sets the given facls on file,
// removing any previous facls. If len(facls)
// == 0, all facls are removed.
func SetFacl(file string, facls ...Facl) error {
	// TODO(synful)
	panic("unimplemented")
}

func faclArgString(facls ...Facl) string {
	var arg string
	for _, facl := range facls {
		arg += facl.aclArg() + ","
	}
	return arg
}

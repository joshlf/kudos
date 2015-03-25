// Package db provides a generic interface to any database backend.
// Packages implementing DBProviders should call RegisterProvider
// in an init function so that it is guaranteed to be available to
// users before their code runs.
package db

import (
	"fmt"
	"sync"
)

type DBKind int

const (
	DBKindString DBKind = iota
	DBKindNumber
)

type DBConstraint struct {
	// TODO(m)
}

type DBEntry struct {
	Kind  DBKind
	Value string
}

type DBEntity interface {
	fmt.Stringer
	RecursiveString() string

	Get(key string) (error, DBEntry)
	Set(key string, value DBEntry)

	AddField(key string, kind DBKind) error
	RemoveField(key string) error

	Children() []DBEntity
	AddChild(name string) error
	RemoveChild(name string) error
}

type DBProvider interface {
	Open() (error, DBProvider)
	Commit() (error, DBProvider)

	Init() (error, DBProvider)
	Destroy() error

	Query(paths []string, constraints []DBConstraint) (error, []DBEntity)
	Modify() error
}

type registry struct {
	providers map[string]DBProvider
	sync.Mutex
}

var reg registry

func init() {
	reg = registry{
		providers: make(map[string]DBProvider),
	}
}

// RegisterProvider registers the given provider under
// the given name. This should only be called by packages
// implementing providers.
func RegisterProvider(name string, provider DBProvider) {
	reg.Lock()
	defer reg.Unlock()
	if _, ok := reg.providers[name]; ok {
		panic(fmt.Sprintf("db: registration of already-registered provider: %v", name))
	}
	reg.providers[name] = provider
}

// GetProvider returns the named provider. It panics
// if the named provider has not previously been
// registered.
func GetProvider(name string) DBProvider {
	reg.Lock()
	defer reg.Unlock()
	p, ok := reg.providers[name]
	if !ok {
		panic(fmt.Sprintf("db: no such provider: %v", name))
	}
	return p
}

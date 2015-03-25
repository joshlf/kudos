// Package db provides a generic interface to any database backend.
// Packages implementing Providers should call RegisterProvider
// in an init function so that it is guaranteed to be available to
// users before their code runs.
package db

import (
	"fmt"
	"sync"
)

type Kind int

const (
	KindString Kind = iota
	KindNumber
)

type Constraint struct {
	// TODO(mdburns)
}

type Entry struct {
	Kind  Kind
	Value string
}

type Entity interface {
	fmt.Stringer
	RecursiveString() string

	Get(key string) (Entry, error)
	Set(key string, value Entry)

	AddField(key string, kind Kind) error
	RemoveField(key string) error

	Children() []Entity
	AddChild(name string) error
	RemoveChild(name string) error
}

type Conn interface {
	Query(path []string, constraints ...Constraint) ([]Entity, error)

	// Close closes the connection and invalidate it;
	// all future calls to Query will fail.
	Close() error
}

type DB interface {
	// Init creates a database if it doesn't
	// already exist. It returns an error
	// if the database already exists.
	Init() error

	// Destroy deletes a database. It returns
	// an error if the database doesn't
	// already exist.
	Destroy() error

	// Connect opens a new connection to
	// a database.
	Connect() (Conn, error)
}

type Provider func(config interface{}) (DB, error)

type registry struct {
	providers map[string]Provider
	sync.Mutex
}

var reg registry

func init() {
	reg = registry{
		providers: make(map[string]Provider),
	}
}

// RegisterProvider registers the given provider under
// the given name. This should only be called by packages
// implementing providers.
func RegisterProvider(name string, provider Provider) {
	reg.Lock()
	defer reg.Unlock()
	if _, ok := reg.providers[name]; ok {
		panic(fmt.Sprintf("db: registration of already-registered provider: %v", name))
	}
	reg.providers[name] = provider
}

// GetDB uses the named provider and given config
// to create a new DB, which is returned. It panics
// if the named provider has not previously been
// registered.
func GetDB(provider string, config interface{}) (DB, error) {
	reg.Lock()
	defer reg.Unlock()
	p, ok := reg.providers[provider]
	if !ok {
		panic(fmt.Sprintf("db: no such provider: %v", provider))
	}
	return p(config)
}

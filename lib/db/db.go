package db

import "fmt"

type DBKind int

const (
	DBKIND_STRING DBKind = iota
	DBKIND_NUMBER
)

type DBEntry struct {
	Kind  DBKind
	Value string
}

type DBEntity interface {
	fmt.Stringer
	RecursiveString() string

	Get(string key) (error, DBEntry)
	Set(string key, value DBEntry)

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

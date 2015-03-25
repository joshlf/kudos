package db

type DBKind int

const (
  DBKIND_STRING DBKind = iota
  DBKIND_NUMBER
)

type DBEntry {
  Kind DBKind
  Value string
}

type DBEntity interface {
  fmt.Stringer
  RecursiveString() string

  Get(string key) (error, DBEntry)
  Set(string key, DBEntry value)

  AddField(string key, DBKind kind) error
  RemoveField(string key) error

  Children() []DBEntity
  AddChild(string name) error
  RemoveChild(string name) error
}

type DBProvider interface {
  Open() (error,DBProvider)
  Commit() (error,DBProvider)

  Init() (error,DBProvider)
  Destroy() error

  Query([]string path, []DBConstraint constraints) (error,[]DBEntity)
  Modify() error
}

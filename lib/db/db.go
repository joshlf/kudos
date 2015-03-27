// Package db provides a generic interface to any database backend.
// Packages implementing Providers should call RegisterProvider
// in an init function so that it is guaranteed to be available to
// users before their code runs.
package db

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/synful/kudos/lib/build"
)

// TODO(synful): here only because it's
// used in the Entity interface; we should
// decide whether it is necessary, and if
// so, use it more consistently
type Kind int

const (
	NumberKind Kind = iota
	StringKind
)

type PathElement struct {
	Any  bool
	Up   bool
	Name string
}

func (p PathElement) String() string {
	if p.Any {
		return "*"
	}
	return p.Name
}

var (
	Any = PathElement{Any: true}
)

type Path struct {
	Path     []PathElement
	Absolute bool
}

func (p *Path) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*p = ParsePath(s)
	return nil
}

func (p Path) String() string {
	var elems []string
	for _, elem := range p.Path {
		elems = append(elems, elem.String())
	}
	str := strings.Join(elems, "/")
	if p.Absolute {
		return "/" + str
	}
	return str
}

// HasWildcards returns whether p
// contains any wildcards
func (p Path) HasWildcards() bool {
	for _, elem := range p.Path {
		if elem.Any {
			return true
		}
	}
	return false
}

// HasUps returns whether p
// contains any .. elements
func (p Path) HasUps() bool {
	for _, elem := range p.Path {
		if elem.Up {
			return true
		}
	}
	return false
}

func ParsePath(path string) Path {
	var p Path
	if len(path) > 0 && path[0] == '/' {
		p.Absolute = true
	}
	fields := strings.FieldsFunc(path, func(r rune) bool { return r == '/' })
	var elem PathElement
	for _, field := range fields {
		switch field {
		case "*":
			elem = PathElement{Any: true}
		case "..":
			elem = PathElement{Up: true}
		default:
			elem = PathElement{Name: field}
		}
		p.Path = append(p.Path, elem)
	}
	return p
}

type Value struct {
	// Value must be string, int64, or float64
	Value interface{}
}

func (v Value) String() string {
	switch v := v.Value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case int64:
		return fmt.Sprint(v)
	case float64:
		// TODO(synful): this distinction
		// may turn out to be unnecessary,
		// but keep it here for now just
		// in case
		if float64(int(v)) == v {
			return fmt.Sprintf("%f.0", v)
		}
		return fmt.Sprint(v)
	}
	if build.DebugMode {
		panic("internal error: Value.Value must be string, int64, or float64")
	}
	return fmt.Sprintf("IllegalValue(%v:%v)", reflect.ValueOf(v.Value), v.Value)
}

func (v *Value) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var n json.Number
		err = json.Unmarshal(b, &n)
		if err != nil {
			return fmt.Errorf("must be string or number")
		}
		i, err := n.Int64()
		if err != nil {
			f, err := n.Float64()
			if err != nil {
				panic(fmt.Errorf("internal error: json.Number is neither float64 nor int64: %v", n))
			}
			v.Value = f
			return nil
		}
		v.Value = i
		return nil
	}
	v.Value = s
	return nil
}

type ConstraintEntity struct {
	Entity interface{}
}

func (c ConstraintEntity) String() string {
	switch e := c.Entity.(type) {
	case Value:
		return fmt.Sprintf("value:%v", e)
	case Path:
		return fmt.Sprintf("path:%v", e)
	}
	if build.DebugMode {
		panic("internal error: ConstraintEntity.Entity must be Value or Path")
	}
	return fmt.Sprintf("UnknownConstraintEntity(%v:%v)", reflect.TypeOf(c.Entity), c.Entity)
}

func (c *ConstraintEntity) UnmarshalJSON(b []byte) error {
	type jsonRepresentation struct {
		Path  *Path
		Value *interface{}
	}
	var j jsonRepresentation
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	if j.Path != nil {
		*c = ConstraintEntity{*j.Path}
	} else {
		*c = ConstraintEntity{*j.Value}
	}
	return nil
}

type Relation int

const (
	LT Relation = iota
	LTE
	GT
	GTE
	EQ
	NEQ
)

var relationKeywords = map[string]Relation{
	"<":  LT,
	"<=": LTE,
	">":  GT,
	">=": GTE,
	"=":  EQ,
	"!=": NEQ,
}

var relationStrings = map[Relation]string{
	LT:  "<",
	LTE: "<=",
	GT:  ">",
	GTE: ">=",
	EQ:  "=",
	NEQ: "!=",
}

func (r Relation) String() string {
	s, ok := relationStrings[r]
	if !ok {
		if build.DebugMode {
			panic(fmt.Errorf("internal error: invalid relation %v", int(r)))
		}
		return fmt.Sprintf("UnknownRelation(%v)", int(r))
	}
	return s
}

func (r *Relation) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	rr, ok := relationKeywords[s]
	if !ok {
		return fmt.Errorf("unknown relation: %v", s)
	}
	*r = rr
	return nil
}

type Constraint struct {
	A, B     ConstraintEntity
	Relation Relation
}

func (c Constraint) String() string {
	return fmt.Sprintf("(%v %v %v)", c.A, c.Relation, c.B)
}

type Entity interface {
	fmt.Stringer
	RecursiveString() string

	Get(key string) (Value, error)
	Set(key string, value Value)

	AddField(key string, kind Kind) error
	RemoveField(key string) error

	Children() []Entity
	AddChild(name string) error
	RemoveChild(name string) error
}

type Conn interface {
	Query(path Path, constraints ...Constraint) ([]Entity, error)

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

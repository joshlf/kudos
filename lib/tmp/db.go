package tmp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func mapToExt(from interface{}, to reflect.Value) error {
	totyp := to.Type()
	fromtyp := typname(reflect.TypeOf(from))

	// NOTE(synful): Numeric conversion code
	// copied from json package's decode.go
	switch totyp.Kind() {
	case reflect.Bool:
		f, ok := from.(bool)
		if !ok {
			return fmt.Errorf("cannot convert %v to %v", reflect.TypeOf(from), totyp)
		}
		to.SetBool(f)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, ok := from.(json.Number)
		if !ok {
			return fmt.Errorf("cannot convert %v to %v", totyp, fromtyp)
		}
		nn, err := strconv.ParseInt(string(n), 10, 64)
		if err != nil || to.OverflowInt(nn) {
			// Since the Number was marshalled by json,
			// we know it's properly formatted, so the
			// only errors ParseInt could return are
			// overflow errors
			return fmt.Errorf("value %v overflows %v", n, totyp)
		}
		to.SetInt(nn)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, ok := from.(json.Number)
		if !ok {
			return fmt.Errorf("cannot convert %v to %v", totyp, fromtyp)
		}
		nn, err := strconv.ParseUint(string(n), 10, 64)
		if err != nil || to.OverflowUint(nn) {
			// Since the Number was marshalled by json,
			// we know it's properly formatted, so the
			// only errors ParseUint could return are
			// overflow errors
			return fmt.Errorf("value %v overflows %v", n, totyp)
		}
		to.SetUint(nn)
	case reflect.Float32, reflect.Float64:
		n, ok := from.(json.Number)
		if !ok {
			return fmt.Errorf("cannot convert %v to %v", totyp, fromtyp)
		}
		nn, err := strconv.ParseFloat(string(n), totyp.Bits())
		if err != nil || to.OverflowFloat(nn) {
			// Since the Number was marshalled by json,
			// we know it's properly formatted, so the
			// only errors ParseFloat could return are
			// overflow errors
			return fmt.Errorf("value %v overflows %v", n, totyp)
		}
		to.SetFloat(nn)
	case reflect.String:
		s, ok := from.(string)
		if !ok {
			return fmt.Errorf("cannot convert %v to %v", totyp, fromtyp)
		}
		to.SetString(s)
	case reflect.Array, reflect.Slice:

	}
	return nil
}

// assumes val.Type() has been validated
func extToMap(val reflect.Value) (interface{}, error) {
	typ := val.Type()
	switch typ.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String,
		reflect.Uintptr, reflect.UnsafePointer:
		// These are all fine to leave as themselves -
		// the json package will marshal them the way
		// we want, so we don't need to do anything custom.
		return val.Interface(), nil
	case reflect.Array, reflect.Slice:
		slc := make([]interface{}, val.Len())
		for i := range slc {
			var err error
			slc[i], err = extToMap(val.Index(i))
			if err != nil {
				return nil, err
			}
		}
		return slc, nil
	case reflect.Interface, reflect.Ptr:
		return extToMap(val.Elem())
	case reflect.Map:
		names := make(map[string]string)
		m := make(map[string]interface{})
		for _, k := range val.MapKeys() {
			effective := sanitizeFieldName(k.String())
			// Can't validate this based on the type
			// alone, so we have to validate dynamically
			if other, ok := names[effective]; ok {
				return nil, fmt.Errorf("map keys %v and %v conflict", k.String(), other)
			}
			names[effective] = k.String()
			var err error
			m[effective], err = extToMap(val.MapIndex(k))
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case reflect.Struct:
		m := make(map[string]interface{})
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			// See https://golang.org/pkg/reflect/#StructField
			exported := field.PkgPath == ""
			if exported {
				var err error
				m[sanitizeFieldName(field.Name)], err = extToMap(val.Field(i))
				if err != nil {
					return nil, err
				}
			}
		}
		return m, nil
	default: // channel, complex, function types
		return nil, fmt.Errorf("unsupported type")
	}
}

func validateType(typ reflect.Type) error {
	return validateTypeHelper(typ, make(map[reflect.Type]bool))
}

func validateTypeHelper(typ reflect.Type, seen map[reflect.Type]bool) error {
	if seen[typ] {
		return nil
	}
	seen[typ] = true

	switch typ.Kind() {
	case reflect.Array, reflect.Slice, reflect.Interface, reflect.Ptr:
		return validateTypeHelper(typ.Elem(), seen)
	case reflect.Map:
		if typ.Key().Kind() != reflect.String {
			return fmt.Errorf("unsupported type %v", typ)
		}
		return validateTypeHelper(typ.Elem(), seen)
	case reflect.Struct:
		m := make(map[string]string)
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			// See https://golang.org/pkg/reflect/#StructField
			exported := field.PkgPath == ""
			if exported {
				effective := sanitizeFieldName(field.Name)
				if other, ok := m[effective]; ok {
					return fmt.Errorf("struct fields %v and %v conflict", field.Name, other)
				}
				m[effective] = field.Name
				if err := validateTypeHelper(field.Type, seen); err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return nil
	}
}

// Removes all but letters and numbers and
// converts to lowercase
func sanitizeFieldName(s string) string {
	s = strings.ToLower(s)
	var ss string
	for _, c := range s {
		switch {
		case 'a' <= c && c <= 'z':
			ss += string(c)
		case '0' <= c && c <= '9':
			ss += string(c)
		}
	}
	return ss
}

// returns "overflows" or "underflows" or the empty string
func overflows(typ reflect.Type, num json.Number) string {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String,
		reflect.Uintptr, reflect.UnsafePointer:
	}
}

var (
	numtype    = reflect.TypeOf(json.Number(""))
	booltype   = reflect.TypeOf(true)
	stringtype = reflect.TypeOf("")
	objtype    = reflect.TypeOf(map[string]interface{}{})
	arrtype    = reflect.TypeOf([]interface{}{})

	typnames = map[reflect.Type]string{
		numtype:    "number",
		booltype:   "bool",
		stringtype: "string",
		objtype:    "object",
		arrtype:    "array",
	}
)

func typname(typ reflect.Type) string {
	name, ok := typnames[typ]
	if !ok {
		panic("internal error: invalid type")
	}
	return name
}

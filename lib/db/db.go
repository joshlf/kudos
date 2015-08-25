package db

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	numtyp    = reflect.TypeOf(json.Number(""))
	booltyp   = reflect.TypeOf(true)
	stringtyp = reflect.TypeOf("")
	objtyp    = reflect.TypeOf(map[string]interface{}{})
	arrtyp    = reflect.TypeOf([]interface{}{})
)

func diff(a, b interface{}, path []interface{}) []interface{} {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return path
	}
	typ := reflect.TypeOf(a)
	switch typ {
	case booltyp, numtyp, stringtyp:
		// TODO(synful): Since we use json.Number
		// for numbers, this amounts to string
		// comparison for numbers. There could be
		// multiple different representations of
		// a given number; we should either convert
		// to an actual number or canonicalize the
		// string before comparing.
		if a != b {
			return path
		}
		return nil
	case objtyp:
		a := a.(map[string]interface{})
		b := b.(map[string]interface{})
		if len(a) != len(b) {
			return path
		}
		for k := range a {
			if _, ok := b[k]; !ok {
				return path
			}
		}
		var changelist [][]interface{}
		for k := range a {
			pathtmp := diff(a[k], b[k], append(path, k))
			if pathtmp != nil {
				changelist = append(changelist, pathtmp)
			}
		}
		switch {
		case len(changelist) == 0:
			return nil
		case len(changelist) == 1:
			return changelist[0]
		default:
			return path
		}
	case arrtyp:
		a := a.([]interface{})
		b := b.([]interface{})
		if len(a) != len(b) {
			return path
		}
		var changelist [][]interface{}
		for i := range a {
			pathtmp := diff(a[i], b[i], append(path, i))
			if pathtmp != nil {
				changelist = append(changelist, pathtmp)
			}
		}
		switch {
		case len(changelist) == 0:
			return nil
		case len(changelist) == 1:
			return changelist[0]
		default:
			return path
		}
	default:
		panic("internal error: unreachable code")
	}
}

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
		f, ok := from.([]interface{})
		if !ok {
			return fmt.Errorf("cannot convert %v to %v", totyp, fromtyp)
		}
		if to.Kind() == reflect.Slice {
			to.Set(reflect.MakeSlice(to.Type(), len(f), len(f)))
		} else if to.Len() != len(f) {
			// TODO(synful): This is as restrictive
			// as can be, so we can always relax the
			// rules later. The json package allows
			// any length, and either doesn't fill
			// the extra etnries in the Go array, or
			// just throws away the extra entries in
			// the json array.
			return fmt.Errorf("cannot convert array of length %v to %v", len(f), fromtyp)
		}
		for i, v := range f {
			err := mapToExt(v, to.Index(i))
			if err != nil {
				return err
			}
		}
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
			if effective == "" {
				return nil, fmt.Errorf("invalid field name %v", k.String())
			}
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
				if effective == "" {
					return fmt.Errorf("invalid struct field name %v", field.Name)
				}
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
// converts to lowercase and removes
// leading numbers
func sanitizeFieldName(s string) string {
	s = strings.ToLower(s)
	var ss string
	var seenLetter bool
	for _, c := range s {
		switch {
		case 'a' <= c && c <= 'z':
			ss += string(c)
			seenLetter = true
		case '0' <= c && c <= '9' && seenLetter:
			ss += string(c)
		}
	}
	return ss
}

var typnames = map[reflect.Type]string{
	numtyp:    "number",
	booltyp:   "bool",
	stringtyp: "string",
	objtyp:    "object",
	arrtyp:    "array",
}

func typname(typ reflect.Type) string {
	name, ok := typnames[typ]
	if !ok {
		panic("internal error: invalid type")
	}
	return name
}

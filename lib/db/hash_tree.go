package db

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

const (
	hashSize = sha256.Size
)

type merkleTree struct {
	hash [hashSize]byte
	// nil for primitive types; key type
	// is string for objects and int for
	// arrays
	children map[interface{}]merkleTree
}

func calcMerkleTree(data interface{}) merkleTree {
	var mt merkleTree
	switch data := data.(type) {
	case json.Number, bool, string:
		// TODO(synful): deal with canonical json.Number issue
		// documented in diff code
		str := fmt.Sprintf("%s/%v", typnames[reflect.TypeOf(data)], data)
		mt.hash = sha256.Sum256([]byte(str))
	case map[string]interface{}:
		// H(obj) = H("object" + H(H(k1)+H(v1)) + H(H(k2)+H(v2)) + ...)
		hash := sha256.New()
		hash.Write([]byte("object"))
		var keys []string
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		mt.children = make(map[interface{}]merkleTree)
		for _, k := range keys {
			// kvhash = H(H(k) + H(v))
			khash := calcMerkleTree(k).hash
			vtree := calcMerkleTree(data[k])
			mt.children[k] = vtree
			kvhash := sha256.Sum256(append(khash[:], vtree.hash[:]...))
			hash.Write(kvhash[:])
		}
		copy(mt.hash[:], hash.Sum(nil))
	case []interface{}:
		// H(arr) = H("array" + H(arr[0]) + H(arr[1]) + ...)
		hash := sha256.New()
		hash.Write([]byte("array"))
		mt.children = make(map[interface{}]merkleTree)
		for i, k := range data {
			t := calcMerkleTree(k)
			mt.children[i] = t
			hash.Write(t.hash[:])
		}
		copy(mt.hash[:], hash.Sum(nil))
	default:
		panic("internal error: unreachable code")
	}
	return mt
}

func merkleDiff(a, b merkleTree, path []interface{}) []interface{} {
	switch {
	case a.hash == b.hash:
		return nil
	case a.children == nil || b.children == nil:
		// One or both are primitive, so
		// a.hash != b.hash means either
		// they're different types or
		// the same primitive type and
		// unequal
		return path
	case len(a.children) != len(b.children):
		return path
	}
	for k := range a.children {
		if _, ok := b.children[k]; !ok {
			return path
		}
	}
	var changelist [][]interface{}
	for k := range a.children {
		pathtmp := merkleDiff(a.children[k], b.children[k], append(path, k))
		if pathtmp != nil {
			changelist = append(changelist, pathtmp)
		}
	}
	if len(changelist) == 1 {
		return changelist[0]
	}
	return path
}

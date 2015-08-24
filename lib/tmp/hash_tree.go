package tmp

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

// This portion of the DB package transforms a JSON blob into a merkle tree.
// Because we intend to store the database in JSON it is possible that someone
// will change the grades database accidentally.  Storing hashes in the
// database will allow us to control this level of human error. It also allows
// for a cleaner interface over which to perform traversals.
// NOTE: there are strong assumptions made here about what type can be put into
// the input of GenHashTree, if anyone changes db.go, the requisite changes
// must be made here as well

// type HashTree struct {
// 	//
// 	Hash     []byte
// 	Payload  interface{}
// 	Children []*HashTree
// 	// this is auxiliary information used to compute the hash of a parent
// 	// It is used for objects, where the string key is used in the hash of the parent
// 	aux []byte
// }

// func hashBaseValue(v interface{}) []byte {
// 	switch v.(type) {
// 	case json.Number, bool, string:
// 		str := fmt.Sprintf("%s/%v", typnames[reflect.TypeOf(v)], v)
// 		sha := sha256.Sum256([]byte(str))
// 		return sha[:]
// 	default:
// 		panic("Invalid input type: this is a programmer error")
// 	}
// }

// Create a hash tree from a json blob (represented as an interface{} with
// proper type assertions. Here we use the sha256 hash. Note the several calls
// to the hash's Write() method without checking the error return value; we do
// this because the method is documented to never return a non-nil error
// func genHashTree(j interface{}) (*HashTree, error) {
// 	res := &HashTree{
// 		Payload: j,
// 	}
// 	switch j.(type) {
// 	case json.Number, bool, string:
// 		res.Hash = hashBaseValue(j)
// 		return res, nil
// 	case map[string]interface{}:
// 		hash := sha256.New()
// 		m := j.(map[string]interface{})
// 		var keys []string
// 		for k := range m {
// 			keys = append(keys, k)
// 		}
// 		sort.Strings(keys)
// 		for _, k := range keys {
// 			ht, err := GenHashTree(m[k])
// 			if err != nil {
// 				//TODO do we want something more informative here? As in
// 				//something that appends to the path that led to the error
// 				return nil, err
// 			}
// 			hash.Write([]byte(k))
// 			hash.Write(ht.Hash)
// 			ht.aux = []byte(k)
// 			res.Children = append(res.Children, ht)
// 		}
// 	case []interface{}:
// 		hash := sha256.New()
// 		for _, k := range j.([]interface{}) {
// 			ht, err := GenHashTree(k)
// 			if err != nil {
// 				return nil, err
// 			}
// 			hash.Write(ht.Hash)
// 			res.Children = append(res.Children, ht)
// 		}
// 	default:
// 		return nil, fmt.Errorf("unsupported type %v", reflect.TypeOf(j))
// 	}
// 	return res, nil
// }

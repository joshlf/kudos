package tmp

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// This portion of the DB package transforms a JSON blob into a merkle tree.
// Because we intend to store the database in JSON it is possible that someone
// will change the grades database accidentally.  Storing hashes in the
// database will allow us to control this level of human error. It also allows
// for a cleaner interface over which to perform traversals.
// NOTE: there are strong assumptions made here about what type can be put into
// the input of GenHashTree, if anyone changes db.go, the requisite changes
// must be made here as well

type HashTree struct {
	//
	Hash     []byte
	Payload  interface{}
	Children []*HashTree
	// this is auxiliary information used to compute the hash of a parent
	// It is used for objects, where the string key is used in the hash of the parent
	aux []byte
}

func hashBaseValue(v interface{}) []byte {
	switch v.(type) {
	case json.Number, bool, string:
		str := fmt.Sprintf("%s/%v", typnames[reflect.TypeOf(v)], v)
		sha := sha256.Sum256([]byte(str))
		return sha[:]
	default:
		panic("Invalid input type: this is a programmer error")
	}
}

// Create a hash tree from a json blob (represented as an interface{} with
// proper type assertions. Here we use the sha256 hash. Note the several calls
// to the hash's Write() method without checking the error return value; we do
// this because the method is documented to never return a non-nil error
func GenHashTree(j interface{}) (*HashTree, error) {
	res := &HashTree{
		Payload: j,
	}
	switch j.(type) {
	case json.Number, bool, string:
		res.Hash = hashBaseValue(j)
		return res, nil
	case map[string]interface{}:
		hash := sha256.New()
		m := j.(map[string]interface{})
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			ht, err := GenHashTree(m[k])
			if err != nil {
				//TODO do we want something more informative here? As in
				//something that appends to the path that led to the error
				return nil, err
			}
			hash.Write([]byte(k))
			hash.Write(ht.Hash)
			ht.aux = []byte(k)
			res.Children = append(res.Children, ht)
		}
	case []interface{}:
		hash := sha256.New()
		for _, k := range j.([]interface{}) {
			ht, err := GenHashTree(k)
			if err != nil {
				return nil, err
			}
			hash.Write(ht.Hash)
			res.Children = append(res.Children, ht)
		}
	default:
		return nil, fmt.Errorf("unsupported type %v", reflect.TypeOf(j))
	}
	return res, nil
}

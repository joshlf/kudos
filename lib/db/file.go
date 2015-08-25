package db

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
)

type dbState struct {
	db    interface{}
	hist  history
	mtree merkleTree
}

type history []change

type change struct {
	Path        path
	Elem        interface{}
	NewRootHash [hashSize]byte
}

type path []interface{}

const hashSize = sha256.Size

type merkleTree struct {
	hash [hashSize]byte
	// nil for primitive types; key type
	// is string for objects and int for
	// arrays
	children map[interface{}]merkleTree
}

func readDB(current, history io.Reader) (dbState, error) {
	cur := json.NewDecoder(current)
	hist := json.NewDecoder(history)
	cur.UseNumber()
	hist.UseNumber()
	var db dbState
	err := cur.Decode(&db.db)
	if err != nil {
		return dbState{}, err
	}
	err = hist.Decode(&db.hist)
	if err != nil {
		return dbState{}, err
	}
	if len(db.hist) == 0 {
		return dbState{}, fmt.Errorf("empty database history")
	}
	db.mtree = calcMerkleTree(db.db)
	if db.mtree.hash != db.hist[0].NewRootHash {
		return dbState{}, fmt.Errorf("current state doesn't match history")
	}
	return db, nil
}

// If old and new are identical, we don't want to clobber the
// old database files. But opening them for writing would
// involve truncating them (doing it otherwise is possible,
// but then figuring out where to truncate after the fact
// is difficult). Thus, the files should only be opened for
// writing (and truncated) when getCur and getHist are called.
func writeDB(getCur, getHist func() (io.Writer, error), old dbState, new interface{}) error {
	newmtree := calcMerkleTree(new)
	p := merkleDiff(old.mtree, newmtree, nil)
	if p == nil {
		return nil
	}
	elem := findElem(new, p)
	change := change{p, elem, newmtree.hash}
	newDBState := dbState{new, append(history(nil), old.hist...), newmtree}
	newDBState.hist = append(newDBState.hist, change)

	current, err := getCur()
	if err != nil {
		return err
	}
	cur := json.NewEncoder(current)
	history, err := getHist()
	if err != nil {
		return err
	}
	hist := json.NewEncoder(history)

	err = cur.Encode(newDBState.db)
	if err != nil {
		return err
	}
	return hist.Encode(newDBState.hist)
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

func merkleDiff(a, b merkleTree, p path) path {
	switch {
	case a.hash == b.hash:
		return nil
	case a.children == nil || b.children == nil:
		// One or both are primitive, so
		// a.hash != b.hash means either
		// they're different types or
		// the same primitive type and
		// unequal
		return p
	case len(a.children) != len(b.children):
		return p
	}
	for k := range a.children {
		if _, ok := b.children[k]; !ok {
			return p
		}
	}
	var changelist []path
	for k := range a.children {
		pathtmp := merkleDiff(a.children[k], b.children[k], append(p, k))
		if pathtmp != nil {
			changelist = append(changelist, pathtmp)
		}
	}
	if len(changelist) == 1 {
		return changelist[0]
	}
	return p
}

// assumes p is in data
func findElem(data interface{}, p path) interface{} {
	switch data.(type) {
	case json.Number, bool, string:
		if len(p) > 0 {
			panic("internal error: findElem called with invalid arguments")
		}
		return data
	case map[string]interface{}:
		data := data.(map[string]interface{})
		return data[p[0].(string)]
	case []interface{}:
		data := data.([]interface{})
		return data[p[0].(int)]
	default:
		panic("internal error: findElem called with invalid arguments")
	}
}

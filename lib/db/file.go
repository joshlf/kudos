package db

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
)

type dbFile struct {
	DB   interface{}
	Hist history
}

type dbState struct {
	dbFile
	mtree merkleTree
}

// changes are sorted in increasing
// order of age (most recent first)
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

func ingestDB(path string) (dbState, error) {
	var db dbState
	var err error
	db.dbFile, err = readDB(path)
	if err != nil {
		return dbState{}, err
	}
	db.mtree = calcMerkleTree(db.DB)
	err = validateDBState(db)
	if err != nil {
		return dbState{}, err
	}
	return db, nil
}

func commitDB(dbpath, tmppath string, old dbState, new interface{}) (bool, error) {
	db, changed := getUpdatedDBState(old, new)
	if !changed {
		return false, nil
	}
	err := writeDB(dbpath, tmppath, db.dbFile)
	if err != nil {
		return false, err
	}
	return true, nil
}

func readDB(path string) (dbFile, error) {
	dbfile, err := os.Open(path)
	if err != nil {
		return dbFile{}, err
	}
	dec := json.NewDecoder(dbfile)
	dec.UseNumber()
	var db dbFile
	err = dec.Decode(&db)
	if err != nil {
		return dbFile{}, err
	}
	return db, nil
}

func writeDB(dbpath, tmppath string, db dbFile) error {
	tmpfile, err := os.Create(tmppath)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(tmpfile)
	err = enc.Encode(db)
	if err != nil {
		return err
	}
	return os.Rename(tmppath, dbpath)
}

func validateDBState(db dbState) error {
	if len(db.Hist) == 0 {
		return fmt.Errorf("empty database history")
	}
	if db.mtree.hash != db.Hist[0].NewRootHash {
		return fmt.Errorf("current state doesn't match history")
	}
	return nil
}

// Computes the new state, including adding
// the diff between old.db and new to the history.
// Returns true if there is a difference, otherwise
// false.
func getUpdatedDBState(old dbState, new interface{}) (dbState, bool) {
	newmtree := calcMerkleTree(new)
	p := merkleDiff(old.mtree, newmtree, nil)
	if p == nil {
		return old, false
	}
	elem := findElem(new, p)
	change := change{p, elem, newmtree.hash}
	newDBState := dbState{dbFile{new, append(history{change}, old.Hist...)}, newmtree}
	return newDBState, true
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

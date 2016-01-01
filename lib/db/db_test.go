package db

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/joshlf/kudos/lib/config"
	"github.com/joshlf/kudos/lib/testutil"
)

/*
	TODO(joshlf):
 		- Test Version and Commit fields
		- Test history functionality
*/

type testDBType struct {
	A uint64
	B float64
}

func randTestDBType() testDBType {
	var t testDBType
	t.A = (uint64(rand.Uint32()) << 32) | uint64(rand.Uint32())
	t.B = rand.Float64()
	return t
}

func TestInit(t *testing.T) {
	tdir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(tdir)
	dbdir := filepath.Join(tdir, "db")
	for i := 0; i < 10; i++ {
		// Initialize a random database,
		// and then read the database back
		// into memory and make sure it
		// matches what we wrote out
		testutil.Must(t, os.Mkdir(dbdir, 0700))
		want := randTestDBType()
		testutil.Must(t, Init(want, dbdir))
		var got testDBType
		_, err := Open(&got, dbdir)
		testutil.Must(t, err)
		testutil.Must(t, os.RemoveAll(dbdir))
		if !reflect.DeepEqual(want, got) {
			t.Errorf("unexpected db value: want %v; got %v", want, got)
		}
	}
}

func TestOpenCommit(t *testing.T) {
	tdir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(tdir)
	want := randTestDBType()
	testutil.Must(t, Init(want, tdir))
	for i := 0; i < 10; i++ {
		// clear to the zero value (in case
		// it is not overwritten and the old
		// value persists)
		got := testDBType{}
		c, err := Open(&got, tdir)
		testutil.Must(t, err)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("unexpected db value: want %v; got %v", want, got)
		}

		// Close the db without committing and make sure
		// the old value of the db persists
		testutil.Must(t, c(nil))
		// clear to the zero value (in case
		// it is not overwritten and the old
		// value persists)
		got = testDBType{}
		c, err = Open(&got, tdir)
		testutil.Must(t, err)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("unexpected db value: want %v; got %v", want, got)
		}

		// Save the new value of the db which we'll
		// check for in the next loop iteration
		want = randTestDBType()
		testutil.Must(t, c(want))
	}
}

func TestLock(t *testing.T) {
	tdir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(tdir)
	var db struct{}
	testutil.Must(t, Init(db, tdir))
	_, err := Open(&db, tdir)
	testutil.Must(t, err)
	_, err = Open(&db, tdir)
	if err != ErrLockFailed {
		t.Errorf("unexpected error: want %v; got %v", err, ErrLockFailed)
	}
}

type marshalError struct {
}

func (m marshalError) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("marshal error")
}

func TestError(t *testing.T) {
	// TODO(joshlf)
	t.Skipf("temporarily skipping until fixed")

	/*
		Open
	*/
	_, err := Open(nil, "")
	testutil.MustError(t, "need absolute path", err)
	_, err = Open(nil, "/dev/null/")
	expect := "acquire lock: open /dev/null/" + config.DBLockFileName + ": not a directory"
	testutil.MustError(t, expect, err)

	tdir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(tdir)
	_, err = Open(nil, tdir)
	expect = "open " + tdir + "/" + config.DBFileName + ": no such file or directory"
	testutil.MustError(t, expect, err)

	f, err := os.Create(filepath.Join(tdir, config.DBFileName))
	testutil.Must(t, err)
	f.Close()

	var db struct{}
	_, err = Open(&db, tdir)
	testutil.MustError(t, "unmarshal from file: EOF", err)

	// make sure that Committer panics if called twice
	testutil.Must(t, Init(db, tdir))
	c, err := Open(&db, tdir)
	testutil.Must(t, err)
	testutil.Must(t, c(db))
	err = func() (err error) {
		defer func() {
			r := recover()
			if r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()
		c(db)
		return nil
	}()
	testutil.MustError(t, "db: Committer called twice", err)

	c, err = Open(&db, tdir)
	testutil.Must(t, err)
	testutil.Must(t, os.Remove(filepath.Join(tdir, config.DBLockFileName)))
	expect = "release lock: remove " + tdir + "/" + config.DBLockFileName +
		": no such file or directory"
	testutil.MustError(t, expect, c(db))

	tdir = testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(tdir)
	testutil.Must(t, Init(db, tdir))
	c, err = Open(&db, tdir)
	testutil.Must(t, err)
	os.Chmod(filepath.Join(tdir, config.DBFileName), 0)
	f, err = os.Create(filepath.Join(tdir, config.DBTempFileName))
	testutil.Must(t, err)
	f.Close()
	testutil.Must(t, os.Chmod(filepath.Join(tdir, config.DBTempFileName), 0))
	expect = "open " + tdir + "/" + config.DBTempFileName + ": permission denied"
	testutil.MustError(t, expect, c(db))

	tdir = testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(tdir)
	testutil.Must(t, Init(db, tdir))
	c, err = Open(&db, tdir)
	testutil.Must(t, err)
	expect = "marshal to file: json: error calling MarshalJSON for type db.marshaler: " +
		"json: error calling MarshalJSON for type db.marshalError: marshal error"
	testutil.MustError(t, expect, c(marshalError{}))

	/*
		Init
	*/
	testutil.MustError(t, "need absolute path", Init(nil, ""))
	testutil.MustError(t, "open /dev/null/"+config.DBFileName+": not a directory", Init(nil, "/dev/null"))

	tdir = testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(tdir)
	expect = "marshal to file: json: error calling MarshalJSON for type db.marshaler: " +
		"json: error calling MarshalJSON for type db.marshalError: marshal error"
	testutil.MustError(t, expect, Init(marshalError{}, tdir))
}

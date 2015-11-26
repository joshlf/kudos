package binexec

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestAssets(t *testing.T) {
	archs := map[string][]struct{ name, sum string }{
		"linux_amd64": {{"tar", "f5b70b4b630f7062b0759c90139322e7e8a27491"}},
		"linux_386":   {{"tar", "98dd383d500cfba6d8adf35d3b9e0c32620c514f"}},
	}

	assets := archs[runtime.GOOS+"_"+runtime.GOARCH]
	for _, a := range assets {
		bin, ok := bins[a.name]
		if !ok {
			t.Errorf("asset not loaded: %v", a.name)
		}
		sum := sha1.Sum(bin.contents)
		str := hex.EncodeToString(sum[:])
		if str != a.sum {
			t.Errorf("asset %v has wrong SHA-1 sum: want %v; got %v", a.name, a.sum, str)
		}
	}
}

const echosrc = `#!/bin/bash

echo -n $1
`

func TestRun(t *testing.T) {
	oldbins := bins
	defer func() {
		bins = oldbins
	}()

	bins = make(map[string]*bin)
	bins["echo"] = &bin{contents: []byte(echosrc)}

	err := Run("nonexistent")
	want := "run \"nonexistent\": no such binary"
	if err == nil || err.Error() != want {
		t.Errorf("unexpected error: want %v; got %v", want, err)
	}

	const rounds = 10
	for i := 0; i < rounds; i++ {
		str := strconv.Itoa(rand.Int())
		var buf bytes.Buffer
		c := exec.Command("echo", str)
		c.Stdout = &buf
		err := RunCmd(c)
		if err != nil {
			t.Errorf("error running subcommand: %v", err)
		}
		out := string(buf.Bytes())
		if out != str {
			t.Errorf("subcommand gave wrong output: want %v; got %v", str, out)
		}
	}

	if bins["echo"].instances != 0 {
		t.Errorf("map entry left with non-zero instances count: %v", bins["echo"].instances)
	}
}

const sleepsrc = `#!/bin/bash

sleep $1
`

func TestConcurrent(t *testing.T) {
	oldbins := bins
	defer func() {
		bins = oldbins
	}()

	bins = make(map[string]*bin)
	bins["sleep"] = &bin{contents: []byte(sleepsrc)}

	const instances = 10
	errs := make([]error, instances)

	var wg sync.WaitGroup
	wg.Add(instances)
	for i := 0; i < 10; i++ {
		go func(i int) {
			errs[i] = Run("sleep", "0.2")
			wg.Done()
		}(i)
	}

	// We want to see where the binary was placed.
	// Sleeping for 0.1s should allow the other
	// goroutines to run, and since they don't
	// sleep or block before acquiring the lock,
	// they should acquire the lock by the time this
	// sleep finishes. Then, they don't release the
	// lock until they've created the binary. However,
	// they WILL block once they're running the
	// subcommand. Thus, we won't acquire the lock
	// until the binary has been written, but almost
	// certainly before all of the subcommands finish,
	// which will allow us to have the lock while the
	// path is still written. Then, later, we can
	// check to make sure it's been properly removed.
	time.Sleep(100 * time.Millisecond)
	binmtx.RLock()
	dirpath := bins["sleep"].dir
	binpath := bins["sleep"].path
	if dirpath == "" {
		t.Log("warning: couldn't get temporary path; " +
			"won't be able to check for proper cleanup")
	}

	fi, err := os.Stat(dirpath)
	if err != nil {
		t.Errorf("could not stat temp directory: %v", err)
	}
	if fi.Mode()&0777 != 0700 {
		t.Errorf("temp directory has wrong perms: want %v; got %v", 07000, fi.Mode()&0777)
	}
	fi, err = os.Stat(binpath)
	if err != nil {
		t.Errorf("could not stat executable: %v", err)
	}
	if fi.Mode()&0777 != 0700 {
		t.Errorf("executable has wrong perms: want %v; got %v", 07000, fi.Mode()&0777)
	}

	binmtx.RUnlock()

	wg.Wait()
	for i, e := range errs {
		if e != nil {
			t.Errorf("unexpected error from instance %v: %v", i, e)
		}
	}

	if dirpath != "" {
		_, err := os.Stat(dirpath)
		if err == nil || !os.IsNotExist(err) {
			t.Errorf("temp directory not cleaned up; stat error was: %v", err)
		}
	}
}

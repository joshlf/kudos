// Package binexec provides the ability to execute binaries
// whose contents have been compiled into the program. This
// is achieved by writing the binaries to a temporary location
// on disk, marking them executable, and executing them.
//
// Currently, the only binaries which are available are those
// which are configured in this package. For Kudos, this is
// sufficient (as far as we've seen), even though it would
// make for bad design if this were a standalone package.
package binexec

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type bin struct {
	contents  []byte
	instances int
	path, dir string
}

var (
	bins   = make(map[string]*bin)
	binmtx sync.RWMutex
)

func writeBin(cmd string) (path string, err error) {
	dirpath, err := ioutil.TempDir("", "kudos")
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dirpath)
		}
	}()

	err = os.Chmod(dirpath, 0700)
	if err != nil {
		return "", err
	}

	bin := bins[cmd]
	binpath := filepath.Join(dirpath, cmd)
	f, err := os.OpenFile(binpath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// umask may (at least on Linux, definitely will)
	// mask out execute bit, so set perms explicitly
	err = os.Chmod(binpath, 0700)
	if err != nil {
		return "", err
	}

	_, err = f.Write(bin.contents)
	if err != nil {
		return "", err
	}

	bins[cmd].dir = dirpath
	bins[cmd].path = binpath
	return binpath, nil
}

// Run runs the binary with the given name and
// arguments. It does this by checking if the
// binary's contents are pre-loaded (compiled
// into this program). If the named binary isn't
// found, an error is returned. Otherwise, the
// binary is written to a temporary location,
// marked executable, and executed with the
// given arguments before being removed.
func Run(name string, args ...string) error {
	return RunCmd(exec.Command(name, args...))
}

// RunCmd is like Run, except that it uses cmd.Args
// as the name of and arguments to the command. It
// is useful if more customization is needed (such
// as redirecting stdin/stdout/stderr).
func RunCmd(cmd *exec.Cmd) error {
	name := cmd.Args[0]

	binmtx.Lock()
	var path string
	bin, ok := bins[name]
	if !ok {
		binmtx.Unlock()
		return fmt.Errorf("run %q: no such binary", name)
	}
	if bin.instances == 0 {
		var err error
		path, err = writeBin(name)
		if err != nil {
			binmtx.Unlock()
			return err
		}
		bin.instances = 1
	} else {
		path = bin.path
	}
	binmtx.Unlock()

	// after this point, bin.instances > 0,
	// so no other instance of Run will delete it

	// we don't want to modify cmd at all, so make
	// a copy of its contents and make sure not to
	// modify any shared memory (namely, the Args
	// slice)
	c := *cmd
	c.Path = path
	c.Args = append([]string(nil), c.Args...)

	err := c.Run()
	binmtx.Lock()
	defer binmtx.Unlock()

	// currently we don't modify the bins map other than
	// to add new entries, but if we ever change that,
	// this would be an easy thing to miss (that is, assuming
	// that the bin pointer is still valid even though we
	// gave up the lock on the mutex for a while), and
	// debugging it would be a pain; better safe than sorry
	bin = bins[name]

	bin.instances--
	if bin.instances == 0 {
		err2 := os.RemoveAll(bin.dir)
		if err2 != nil && err == nil {
			err = err2
		}
		// allow gc; minor, but no reason not to
		bin.path = ""
		bin.dir = ""
	}
	return err
}

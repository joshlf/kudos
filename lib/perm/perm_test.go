package perm

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	rand.Seed(30682)
	on := "rwxrwxrwx"
	off := "---------"
	for i := 0; i < 1000; i++ {
		var str string
		for i, c := range on {
			if rand.Int()%2 == 0 {
				c = rune(off[i])
			}
			str += string(c)
		}
		got := Parse(str).String()[1:]
		if got != str {
			t.Errorf("parsed %v, got %v", str, got)
		}
	}
}

func TestAddFacls(t *testing.T) {
	// There's no good way to test that the facls
	// work since you would need to not be the
	// primary user on the file in order for the
	// user facls to take precedence.  Thus, just
	// test whether the facls have been applied;
	// don't test whether they work.

	u, err := user.Current()
	if err != nil {
		t.Fatalf("could not get current user: %v", err)
	}
	out, err := exec.Command("getent", "group", fmt.Sprint(u.Gid)).CombinedOutput()
	if err != nil {
		t.Fatalf("could not get user's primary group name: %v", err)
	}
	fields := strings.FieldsFunc(string(out), func(r rune) bool { return r == ':' || r == '\n' })
	if len(fields) != 3 || fields[0] == "" {
		t.Fatalf("bad output from `getent group %v`: %v", u.Gid, string(out))
	}
	gname := fields[0]

	const fname = "test_set_facl"

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_EXCL, 0)
	if err != nil {
		t.Fatalf("could not create test file: %v", err)
	}
	f.Close()
	defer os.Remove(fname)

	err = AddFacl(fname, Facl{User, u.Username, Read}, Facl{Group, u.Gid, Write}, Facl{Entity: Other, Perm: Execute})
	if err != nil {
		t.Fatalf("could not add facls: %v", err)
	}

	out, err = exec.Command("getfacl", fname).CombinedOutput()
	if err != nil {
		t.Fatalf("could not get facls: %v", err)
	}
	lines := strings.FieldsFunc(string(out), func(r rune) bool { return r == '\n' })
	expect := []string{
		fmt.Sprintf("# file: %v", fname),
		fmt.Sprintf("# owner: %v", u.Username),
		fmt.Sprintf("# group: %v", u.Username),
		fmt.Sprintf("user::---"),
		fmt.Sprintf("user:%v:r--", u.Username),
		fmt.Sprintf("group::---"),
		fmt.Sprintf("group:%v:-w-", gname),
		fmt.Sprintf("mask::rw-"),
		fmt.Sprintf("other::--x"),
	}
	if !reflect.DeepEqual(lines, expect) {
		t.Errorf("unexpected output from getfacl; want:\n%v\n\ngot:\n%v", strings.Join(expect, "\n"), strings.Join(lines, "\n"))
	}
}

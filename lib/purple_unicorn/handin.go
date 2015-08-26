package purple_unicorn

// This file provides functions to manage the public handin metadata file for a
// course. This construct creates a public manifest that is readable and
// writeable by TAs, and readable by students

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/synful/go-acl"
)

//name of file where where public handin specs are stored.
const HandinSpecFile = "pubfile"

type PublicHandins []PublicHandinSpec

type PublicHandinSpec struct {
	Assignment string
	Files      []Code
	DueDate    time.Time // how are we representing Dates?
}

// represents a 'path' starting from the root of a course category to either an
// Assignment or a handin, whichever is most specific
type AssignmentPath []string

func (c *Course) PubFile() string {
	return path.Join(c.location, HandinSpecFile)
}

// This file should create the public file with ACL read/writes to the TA group
// and acl reads to the student group
func (c *Course) CreatePubFile() error {
	_, err := os.Create(c.PubFile())
	if err != nil {
		return err
	}
	acl.Set(c.PubFile(), acl.ACL([]acl.Entry{acl.Entry{
		Tag:       acl.TagGroup,
		Qualifier: string(*c.studentGroup),
		Perms:     4,
	}, acl.Entry{
		Tag:       acl.TagGroup,
		Qualifier: string(*c.taGroup),
		Perms:     6,
	}}))

	return nil
}

func (c *Course) AppendSpec(p PublicHandins) error {
	var spec PublicHandins
	f, err := os.OpenFile(c.PubFile(), os.O_RDWR, os.FileMode(0))
	if err != nil {
		return err
	}
	d := json.NewDecoder(f)
	e := json.NewEncoder(f)
	err = d.Decode(&spec)
	if err != nil {
		return err
	}
	spec = append(spec, p...)
	return e.Encode(&spec)
}
func (c *Course) WriteSpec(p PublicHandins) error {
	f, err := os.OpenFile(c.PubFile(), os.O_WRONLY, os.FileMode(0))
	if err != nil {
		return err
	}
	e := json.NewEncoder(f)
	return e.Encode(&p)
}

func (c *Course) RemoveSpec(p PublicHandinSpec) error {
	var spec PublicHandins
	f, err := os.OpenFile(c.PubFile(), os.O_RDWR, os.FileMode(0))
	if err != nil {
		return err
	}
	d := json.NewDecoder(f)
	e := json.NewEncoder(f)
	err = d.Decode(&spec)
	if err != nil {
		return err
	}
	found := false
	ind := 0
	for i, hspec := range spec {
		if hspec.Assignment == p.Assignment && reflect.DeepEqual(hspec.Files, p.Files) {
			found, ind = true, i
			break
		}
	}
	if !found {
		return fmt.Errorf("cannot find handin")
	}
	spec = append(spec[:ind], spec[ind+1:]...)
	return e.Encode(&spec)
}

func (c *Category) OpenHandin(path AssignmentPath) ([]PublicHandinSpec, error) {
	res := []PublicHandinSpec{}
	cur := c
	var asgn *Assignment = nil
	for _, str := range path {
		if len(cur.Assignments()) == 0 {
			found := false
			for _, child := range cur.Children() {
				if child.name == str {
					found = true
					cur = child
					break
				}
			}
			if !found {
				return res, fmt.Errorf("unable to find subcategory or assignment %v in category %v", str, c.Name())
			}
		} else if asgn == nil {
			// We have reached the end of the category hierarchy, now we need to find an assignment
			found := false
			for _, asg := range cur.Assignments() {
				//TODO should code be used here?
				if asg.EffectiveName() == str {
					found = true
					asgn = asg
					break
				}
			}
			if !found {
				return res, fmt.Errorf("unable to find subcategory or assignment %v in category %v", str, c.Name())
			}
		} else {
			// we have found an assignment, now to find a handin
			for _, handin := range asgn.Handins() {
				if handin.hasCode() && handin.Code == Code(str) {
					return []PublicHandinSpec{PublicHandinSpec{
						Assignment: asgn.EffectiveName(),
						Files:      handin.problems(),
						DueDate:    handin.due(),
					}}, nil
				}
			}
		}
	}
	// The entire assignment is being handed in.
	for _, handin := range asgn.Handins() {
		res = append(res, PublicHandinSpec{
			Assignment: asgn.EffectiveName(),
			Files:      handin.problems(),
			DueDate:    handin.due(),
		})
	}

	return res, nil

}

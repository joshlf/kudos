package purple_unicorn

import (
	"fmt"
	"io/ioutil"
	"os"
)

const CategoryConfigFileName = "config.json"

//This file contains experimental work on dealing with course 'categories': A
//layer of abstraction that sits on top of assignments to assign weights
//between groups of assignments. I.e. A course may have different assignments
//be designated as "Homework" or "Projects". This framework allows arbitrary
//trees of categories.
//From a UI/organizational perspective, we have two options on how to organize
//categories:
// - Force each config file to list all of its children explicitly
//    (e.g. "children"  : ["foo/file1.json", "bar/file2.json"...])
// - Deduce the location of subcategories based on subdirectories.
//Luckily the overall API will remain the same regardless, but this does
//significantly change parsing

type Weight uint64

//Recursively traverse directory tree, building up category structure.
//The base case of the recursion is the root category generated from a course config along
//with the directory to start searching for configs
func ParseCategory(parent *Category, dir string) error {
	var (
		errs       ErrList
		toTraverse []os.FileInfo
		newParent  *Category
	)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory %v: %v", dir, err)
	}

start:
	for _, f := range files {
		fileName := fmt.Sprintf("%s/%s", dir, f.Name())
		if f.IsDir() {
			if newParent == nil {
				// We have not found a course config file, so set aside this
				// directory to process later on
				toTraverse = append(toTraverse, f)
			} else {
				err = ParseCategory(newParent, fileName)
				if err != nil {
					errs.Add(fmt.Errorf("unable to parse directory %s/%s: %v", dir, f.Name(), err))
				}
			}
		} else if f.Name() == CategoryConfigFileName {
			// Found a category config file
			r, err := os.Open(fileName)
			if err != nil {
				errs.Add(fmt.Errorf("unable to open file %s/%s", dir, f.Name()))
				// return now, because later errors will not make very much sense
				// if there is no config.
				return errs
			}
			if b, err := parseCategory(r); err != nil {
				errs.Add(fmt.Errorf("error parsing category config %s/%s: %v", dir, f.Name(), err))
				return errs
			} else {
				if newParent != nil {
					// this is a programmer error as it would imply that there
					// is more than one file in this directory with this name
					panic("Internal Error: Multiple config files found")
				}
				newParent = &Category{
					name:   *b.Name,
					weight: Weight(*b.Weight),
				}
				parent.children = append(parent.children, newParent)
			}
		} else {
			// We have found a file -- presumably this means it is an assignment
			a, err := ParseAssignmentFile(fileName)
			if err != nil {
				errs.Add(err)
			} else {
				parent.asgns = append(parent.asgns, a)
			}
		}
	}

	//If we encountered some directories before parsing the config file, then
	//we need to traverse them
	if len(toTraverse) != 0 {
		files = toTraverse
		toTraverse = nil
		goto start
	}

	return errs
}

type Category struct {
	name     string
	children []*Category
	asgns    []*Assignment
	weight   Weight
}

func (b *Category) Validate() error {
	if b.children != nil && b.asgns != nil {
		return fmt.Errorf(`Category %v has the following children that are subcatecories:
%v
and the following children that are assignments:
%v
Categories may only have one or the other`, b.name, b.children, b.asgns)
	}
	return nil
}

func (b *Category) Children() []*Category {
	return b.children
}

func (b *Category) Assignments() []*Assignment {
	return b.asgns
}

func (b *Category) Weight() Weight {
	return b.weight
}

func (b *Category) Name() string {
	return b.name
}

// A similar recursive strategy can/should be used to generate a
// pretty-printing method for showing what the course structure looks like.  A
// more general framework along those lines would be to allow for it to render
// html/latex for course documentation.
func GetAssignments(c *Category) []*Assignment {
	var res []*Assignment
	if children := c.Children(); len(children) != 0 {
		for _, ch := range children {
			res = append(res, GetAssignments(ch)...)
		}
	} else {
		res = c.Assignments()
	}
	return res
}

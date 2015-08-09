package purple_unicorn

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Group string
type User string

type Course struct {
	code         Code
	name         string
	description  string
	taGroup      *Group
	tas          []User //TODO:  []TA?, how should this be represented
	studentGroup *Group
	handinMethod string
}

func (g Group) Validate() error {
	group := string(g)
	if f, err := os.Open("/etc/group"); err == nil {
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("Could not verify group: %v\n", buf)
		}
		str := string(buf)
		for _, line := range strings.Split(str, "\n") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 && parts[0] == group {
				return nil
			}
		}
		return fmt.Errorf("no such group: %v", group)
	} else {
		//TODO have a disciplined way of validating groups.
		return nil
	}
}

//TODO: we probably want to have separate Group and TA types to do a validate on
func (u User) Validate() error {
	//TODO nyi
	return nil
}

func (c *Course) Validate() error {
	var errs ErrList

	if err := c.code.Validate(); err != nil {
		errs = append(errs, err)
	}
	if c.taGroup != nil {
		if len(*c.taGroup) == 0 {
			errs = append(errs,
				fmt.Errorf("Ta Group must not be the empty string; if you don't want a ta group, remove the line"))
			goto students
		}
		if err := c.taGroup.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("Ta group invalid: %v", err))
		}
	}

students:
	if c.studentGroup != nil {
		if len(*c.studentGroup) == 0 {
			errs = append(errs, fmt.Errorf("Student Group must not be the empty string; if you don't want a ta group, remove the line"))
			goto studentList
		}
		if err := c.studentGroup.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("Student group invalid: %v", err))
		}
	}

studentList:
	var localErrList ErrList
	for _, ta := range c.tas {
		if err := ta.Validate(); err != nil {
			localErrList = append(errs, err)
		}
	}
	if len(localErrList) > 0 {
		errs = append(errs, localErrList)
	}
	if len(errs) > 0 {
		return fmt.Errorf("Error validating Course config:\n%v", errs)
	}

	return nil
}

package purple_unicorn

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type course struct {
	Code         *string  `json:"code"`
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	TaGroup      *string  `json:"ta_group"`
	TAs          []string `json:"tas"`
	StudentGroup *string  `json:"student_group"`
	HandinMethod *string  `json:"handin_method"`
}

func ParseCourseFile(f string) (*Course, error) {
	if r, err := os.Open(f); err == nil {
		if p, err := ParseCourse(r); err != nil {
			return nil, fmt.Errorf("error parsing file %v:\n%v", f, err)
		} else {
			return p, nil
		}
	} else {
		return nil, fmt.Errorf("unable to open course file %v", f)
	}
}

func ParseCourse(r io.Reader) (*Course, error) {
	var errs ErrList
	d := json.NewDecoder(r)
	var c course
	err := d.Decode(&c)
	if err != nil {
		return nil, err
	}
	course := Course{}
	if c.Code == nil {
		errs.Add(fmt.Errorf("must add \"code\" field"))
	} else {
		course.SetCodeNoValidate(Code(*c.Code))
	}
	if c.Name == nil {
		errs.Add(fmt.Errorf("must add \"name\" field"))
	} else {
		course.SetNameNoValidate(*c.Name)
	}
	if c.Description != nil {
		course.SetDescriptionNoValidate(*c.Description)
	}
	if c.TaGroup != nil {
		taGroup := Group(*c.TaGroup)
		course.SetTaGroupNoValidate(&taGroup)
	}
	for _, ta := range c.TAs {
		course.AddTANoValidate(User(ta))
	}
	if c.StudentGroup != nil {
		studGroup := Group(*c.StudentGroup)
		course.SetStudentGroupNoValidate(&studGroup)
	}
	course.SetHandinMethodNoValidate(c.HandinMethod)

	if len(errs) > 0 {
		return nil, errs
	}
	return &course, course.Validate()
}

func NewCourse(code Code, name string, description string, taGroup *Group,
	tas []User, studentGroup *Group, handinMethod *string) Course {
	return Course{code, name, description, taGroup, tas, studentGroup, handinMethod}
}

type Course struct {
	code         Code
	name         string
	description  string
	taGroup      *Group
	tas          []User
	studentGroup *Group
	handinMethod *string
}

// supported handin methods
var handinMethods = map[string]bool{
	"facl": true,
}

type Group string

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

func (g Group) MustValidate() {
	err := g.Validate()
	if err != nil {
		panic(err)
	}
}

type User string

//TODO: we probably want to have separate Group and TA types to do a validate on
func (u User) Validate() error {
	//TODO nyi
	return nil
}

func (u User) MustValidate() {
	err := u.Validate()
	if err != nil {
		panic(err)
	}
}

func (c *Course) Validate() error {
	var errs ErrList

	if err := c.code.Validate(); err != nil {
		errs.Add(err)
	}

	if c.taGroup != nil {
		if len(*c.taGroup) == 0 {
			errs.Add(fmt.Errorf("Ta Group must not be the empty string; if you don't want a ta group, remove the line"))
			goto students
		}
		if err := c.taGroup.Validate(); err != nil {
			errs.Add(fmt.Errorf("Ta group invalid: %v", err))
		}
	}

students:
	if c.studentGroup != nil {
		if len(*c.studentGroup) == 0 {
			errs.Add(fmt.Errorf("Student Group must not be the empty string; if you don't want a ta group, remove the line"))
			goto studentList
		}
		if err := c.studentGroup.Validate(); err != nil {
			errs.Add(fmt.Errorf("Student group invalid: %v", err))
		}
	}

studentList:
	var localErrList ErrList
	for _, ta := range c.tas {
		if err := ta.Validate(); err != nil {
			localErrList.Add(err)
		}
	}
	if len(localErrList) > 0 {
		errs.Add(localErrList)
	}
	if c.handinMethod != nil && !handinMethods[*c.handinMethod] {
		errs.Add(fmt.Errorf("Unsupported handin method: %v", *c.handinMethod))
	}

	if len(errs) > 0 {
		return fmt.Errorf("Error validating Course config:\n%v", errs)
	}

	return nil
}

func (c *Course) MustValidate() {
	err := c.Validate()
	if err != nil {
		panic(err)
	}
}

func (c *Course) SetTaGroup(ta *Group) (old *Group, err error) {
	oTa := c.taGroup
	c.SetTaGroupNoValidate(ta)
	return oTa, c.Validate()
}

func (c *Course) SetTaGroupNoValidate(ta *Group) {
	c.taGroup = ta
}

func (c *Course) SetStudentGroup(s *Group) (old *Group, err error) {
	o := c.studentGroup
	c.SetStudentGroupNoValidate(s)
	return o, c.Validate()
}

func (c *Course) SetStudentGroupNoValidate(student *Group) {
	c.studentGroup = student
}
func (c *Course) SetCode(co Code) (old Code, err error) {
	oldCode := c.code
	c.SetCodeNoValidate(co)
	return oldCode, c.Validate()
}

func (c *Course) SetCodeNoValidate(co Code) {
	c.code = co
}

func (c *Course) SetName(n string) (old string, err error) {
	oldName := c.name
	c.SetNameNoValidate(n)
	return oldName, c.Validate()
}

func (c *Course) SetNameNoValidate(n string) {
	c.name = n
}

func (c *Course) SetDescription(d string) (old string, err error) {
	oldD := c.description
	c.SetDescriptionNoValidate(d)
	return oldD, c.Validate()
}

func (c *Course) SetDescriptionNoValidate(d string) {
	c.description = d
}

func (c *Course) AddTA(ta User) error {
	c.AddTANoValidate(ta)
	if err := c.Validate(); err != nil {
		c.tas = c.tas[:len(c.tas)-2]
		return err
	}
	return nil
}

func (c *Course) AddTANoValidate(ta User) {
	c.tas = append(c.tas, ta)
}

func (c *Course) SetHandinMethod(h *string) (old *string, err error) {
	o := c.handinMethod
	c.SetHandinMethodNoValidate(h)
	return o, c.Validate()
}

func (c *Course) SetHandinMethodNoValidate(h *string) {
	c.handinMethod = h
}

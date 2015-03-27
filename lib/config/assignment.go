package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

func timeparse(text string) (time.Time, error) {
	return time.Parse("Jan 2, 2006 at 3:04pm (MST)", text)
}

type GradeNum float64

func (g *GradeNum) UnmarshalTOML(i interface{}) error {
	f, ok := i.(float64)
	if !ok {
		ii, ok := i.(int64)
		if !ok {
			return fmt.Errorf("expected number")
		}
		f = float64(ii)
	}
	*g = GradeNum(f)
	return nil
}

func (g GradeNum) MarshalText() ([]byte, error) {
	if GradeNum(int(g)) == g {
		return []byte(fmt.Sprint(int(g))), nil
	}
	return []byte(fmt.Sprint(g)), nil
}

type Problem struct {
	Name  string
	Files []string
	Total GradeNum
}

func DefaultProblem() Problem {
	return Problem{
		Name:  "Problem Name",
		Files: []string{"file1.txt"},
		Total: GradeNum(100),
	}
}

type AssignSpec struct {
	// TODO(jliebowf): require assignment names to be
	// sanitized (e.g., contain no spaces)? would allow
	// us to assume we can use them as handin dir names
	Name    string
	Problem []Problem
	Handin  Handin
}

func DefaultAssignSpec() AssignSpec {
	return AssignSpec{
		Name:    "Assignment Name",
		Problem: []Problem{DefaultProblem()},
		Handin:  Handin{},
	}
}

type Handin struct {
	Due   date
	Grace duration
}

type date struct {
	time.Time
}

func (d *date) UnmarshalText(text []byte) error {
	var err error
	d.Time, err = timeparse(string(text))
	return err
}

// TAKEN from burntsushi example
type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

//end taken from example

func AsgnFromFile(file string) (AssignSpec, error) {
	var fileText []byte
	var err error
	var res AssignSpec

	if fileText, err = ioutil.ReadFile(file); err != nil {
		return AssignSpec{}, err
	}

	if _, err = toml.Decode(string(fileText), &res); err != nil {
		return AssignSpec{}, err
	}

	return res, nil
}

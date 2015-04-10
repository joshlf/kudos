package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

func ReadAllAssignments(course Course) ([]Assignment, error) {
	adir := course.AssignmentsDir()
	entries, err := ioutil.ReadDir(adir)
	if err != nil {
		return nil, err
	}
	var a []Assignment
	for _, e := range entries {
		// TODO(synful): e.IsDir() isn't really good enough,
		// but checking e.Mode().IsRegular() doesn't work
		// either because we want to support symlinks, for
		// example
		name := e.Name()
		if !e.IsDir() && len(name) > 5 && name[len(name)-5:] == ".toml" {
			path := filepath.Join(adir, name)
			conf, err := readAssignConfig(name[:len(name)-5], path)
			if err != nil {
				return nil, err
			}
			a = append(a, Assignment{conf, course})
		}
	}
	return a, nil
}

func ReadAssignment(course Course, code string) (Assignment, error) {
	adir := course.AssignmentsDir()
	// TODO(synful): do we want to assume toml?
	file := filepath.Join(adir, code+".toml")
	conf, err := readAssignConfig(code, file)
	if err != nil {
		return Assignment{}, err
	}
	return Assignment{conf, course}, nil
}

type Assignment struct {
	conf   assignConfig
	course Course
}

// Code returns a's code.
func (a Assignment) Code() string { return string(a.conf.Code.code) }

// Name returns the human-readbale
// name of a. If one was not set
// in the config file, it defaults
// to a.Code().
func (a Assignment) Name() string {
	if !a.conf.Name.set {
		return a.Code()
	}
	return a.conf.Name.string
}

// HasMultipleHandins returns whether a
// has multiple handins.
func (a Assignment) HasMultipleHandins() bool { return !a.conf.Due.set }

// Handin returns all of the handins for a.
// If a.HasMultipleHandins() == false, it
// panics; in this case, callers should
// instead use a.Due().
func (a Assignment) Handins() []Handin {
	if !a.HasMultipleHandins() {
		panic("config: does not have multiple handins")
	}
	var h []Handin
	for _, hh := range a.conf.Handins {
		h = append(h, hh.toHandin(a.conf.Problems))
	}
	return h
}

// Handin returns the due date for a.
// If a.HasMultipleHandins() == true,
// it panics; in this case, callers
// should instead use a.Handins().
func (a Assignment) Due() time.Time {
	if a.HasMultipleHandins() {
		panic("config: has multiple handins")
	}
	return time.Time(a.conf.Due.date)
}

func (a Assignment) HandinDir() string {
	return filepath.Join(a.course.HandinDir(), string(a.conf.Code.code))
}

func (a Assignment) Problems() []Problem {
	var p []Problem
	for _, pp := range a.conf.Problems {
		p = append(p, pp.toProblem())
	}
	return p
}

type Problem struct {
	Code        string
	name        optionalString
	points      float64
	SubProblems []Problem
}

// Name returns the human-readbale
// name of p. If one was not set
// in the config file, it defaults
// to p.Code.
func (p Problem) Name() string {
	if !p.name.set {
		return p.Code
	}
	return p.name.string
}

// Points returns the number of points that
// p is worth. If len(p.SubProblems) > 0,
// it will be inferred from the point values
// of its subproblems. Otherwise, it will
// will be the specified point value of
// the problem itself. It is guaranteed that
// this package will never return a Problem
// which has either neither points nor
// subproblems, or which has both.
func (p Problem) Points() float64 {
	if len(p.SubProblems) == 0 {
		return p.points
	}
	points := float64(0)
	for _, pp := range p.SubProblems {
		points += pp.Points()
	}
	return points
}

type Handin struct {
	Code     string
	Due      time.Time
	Problems []Problem
}

type assignConfig struct {
	// Guranteed to be set
	Code     optionalCode   `toml:"code"`
	Name     optionalString `toml:"name"`
	Due      optionalDate   `toml:"due"` // set if and only if len(Handins) == 0
	Handins  []handin       `toml:"handin"`
	Problems []problem      `toml:"problem"`
}

type problem struct {
	// Guaranteed to be set
	Code optionalCode   `toml:"code"`
	Name optionalString `toml:"name"`

	// Guaranteed that one of these fields
	// will be set, but not both.
	Points      optionalNumber `toml:"points"`
	SubProblems []problem      `toml:"subproblem"`
}

func (p problem) toProblem() Problem {
	pp := Problem{
		Code:   string(p.Code.code),
		name:   p.Name,
		points: float64(p.Points.number),
	}
	for _, ppp := range p.SubProblems {
		pp.SubProblems = append(pp.SubProblems, ppp.toProblem())
	}
	return pp
}

type handin struct {
	// All fields are guaranteed to be set
	Code     optionalCode `toml:"code"`
	Due      optionalDate `toml:"due"`
	Problems []code       `toml:"problems"`
}

func (h handin) toHandin(p []problem) Handin {
	probs := make(map[code]Problem)
	for _, pp := range p {
		probs[pp.Code.code] = pp.toProblem()
	}

	hh := Handin{
		Code: string(h.Code.code),
		Due:  time.Time(h.Due.date),
	}
	for _, p := range h.Problems {
		pp, ok := probs[p]
		if !ok {
			panic("config: internal error")
		}
		hh.Problems = append(hh.Problems, pp)
	}
	return hh
}

func readAssignConfig(code, file string) (assignConfig, error) {
	var conf assignConfig
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return assignConfig{}, err
	}
	if !conf.Code.set {
		return assignConfig{}, fmt.Errorf("assignment must have code")
	}
	if string(conf.Code.code) != code {
		return assignConfig{}, fmt.Errorf("assignment code in config (%v) does not match expected code (%v)", string(conf.Code.code), code)
	}
	if len(conf.Problems) == 0 {
		return assignConfig{}, fmt.Errorf("assignment has no problems")
	}
	if conf.Due.set && len(conf.Handins) != 0 {
		return assignConfig{}, fmt.Errorf("assignment cannot have due date and handins")
	}
	if len(conf.Handins) == 1 {
		return assignConfig{}, fmt.Errorf("assignment cannot have one handin - instead just use a due date")
	}

	var validateProblem func(path string, p problem) error
	validateProblem = func(path string, p problem) error {
		if !p.Code.set {
			// Don't print the path since there was no code,
			// so the path won't identify the problem.
			return fmt.Errorf("all problems must have a code")
		}
		if p.Points.set && len(p.SubProblems) > 0 {
			return fmt.Errorf("problem cannot have points and subproblems: %v", path)
		}
		if !p.Points.set && len(p.SubProblems) == 0 {
			return fmt.Errorf("problem must have points or subproblems: %v", path)
		}
		probs := make(map[string]bool)
		for _, pp := range p.SubProblems {
			if probs[string(pp.Code.code)] {
				return fmt.Errorf("duplicate subproblem code: %v.%v", path, pp.Code.code)
			}
			probs[string(pp.Code.code)] = true
			if err := validateProblem(path+"."+string(pp.Code.code), pp); err != nil {
				return err
			}
		}
		return nil
	}

	probs := make(map[string]bool)
	for _, p := range conf.Problems {
		if probs[string(p.Code.code)] {
			return assignConfig{}, fmt.Errorf("duplicate problem code: %v", p.Code.code)
		}
		probs[string(p.Code.code)] = true
		if err := validateProblem(string(p.Code.code), p); err != nil {
			return assignConfig{}, err
		}
	}
	// now probs contains all problems
	if !conf.Due.set {
		for _, h := range conf.Handins {
			if !h.Code.set {
				return assignConfig{}, fmt.Errorf("all handins must have a code")
			}
			if !h.Due.set {
				return assignConfig{}, fmt.Errorf("handin must have due date: %v", h.Code.code)
			}
			if len(h.Problems) == 0 {
				return assignConfig{}, fmt.Errorf("handin must have problems: %v", h.Code.code)
			}
			for _, p := range h.Problems {
				notSeen, ok := probs[string(p)]
				if !ok {
					return assignConfig{}, fmt.Errorf("unknown problem in handin %v: %v", h.Code.code, p)
				}
				if !notSeen {
					return assignConfig{}, fmt.Errorf("problem in multiple handins: %v", p)
				}
				probs[string(p)] = false
			}
		}
		for p, notSeen := range probs {
			if notSeen {
				return assignConfig{}, fmt.Errorf("problem not in any handins: %v", p)
			}
		}
	}
	return conf, nil
}

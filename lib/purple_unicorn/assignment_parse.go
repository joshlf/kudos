package purple_unicorn

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type handin struct {
	Code     *string  `json:"code"`
	Due      *date    `json:"due"`
	Problems []string `json:"problems"`
}

// Convert h to an exported Handin type.
// This function performs no validation,
// so you must do validation independent
// of this function.
func (h handin) toHandin() (hh Handin) {
	hh.Code = h.code()
	hh.HasCode = h.hasCode()
	hh.Due = h.due()
	hh.Problems = h.problems()
	return
}

// handin implements the parseableHandinInterface interface
func (h handin) code() (c Code) {
	if h.Code != nil {
		c = Code(*h.Code)
	}
	return
}

func (h handin) due() (t time.Time) {
	if h.Due != nil {
		t = time.Time(*h.Due)
	}
	return
}

func (h handin) problems() (probs []Code) {
	for _, pp := range h.Problems {
		probs = append(probs, Code(pp))
	}
	return
}

func (h handin) hasCode() bool { return h.Code != nil }
func (h handin) hasDue() bool  { return h.Due != nil }

type problem struct {
	Code        *string   `json:"code"`
	Name        *string   `json:"name"`
	Points      *float64  `json:"points"`
	Subproblems []problem `json:"subproblems"`
}

// Convert p to an exported Problem type.
// This function performs no validation,
// so you must do validation independent
// of this function.
func (p problem) toProblem() (pp Problem) {
	pp.Code = p.code()
	pp.Name = p.name()
	pp.Points = p.points()
	pp.HasPoints = p.points == nil
	for _, ppp := range p.Subproblems {
		pp.Subproblems = append(pp.Subproblems, ppp.toProblem())
	}
	return
}

// problem implements the parseableProblemInterface interface
func (p problem) code() (c Code) {
	if p.Code != nil {
		c = Code(*p.Code)
	}
	return
}

func (p problem) name() (s string) {
	if p.Name != nil {
		s = *p.Name
	}
	return
}

func (p problem) points() (f float64) {
	if p.Points != nil {
		f = *p.Points
	}
	return
}

func (p problem) subproblems() (probs []problemInterface) {
	for _, pp := range p.Subproblems {
		probs = append(probs, pp)
	}
	return
}

func (p problem) hasCode() bool   { return p.Code != nil }
func (p problem) hasName() bool   { return p.Name != nil }
func (p problem) hasPoints() bool { return p.Points != nil }

type assignment struct {
	Code     *string   `json:"code"`
	Name     *string   `json:"name"`
	Handins  []handin  `json:"handins"`
	Problems []problem `json:"problems"`
}

func ParseAssignmentFile(path string) (*Assignment, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parse %v: %v", path, err)
	}
	a, err := parseAssignment(f)
	if err != nil {
		return nil, fmt.Errorf("parse %v: %v", path, err)
	}
	return a, nil
}

func ParseAssignment(r io.Reader) (*Assignment, error) {
	a, err := parseAssignment(r)
	if err != nil {
		return nil, fmt.Errorf("parse: %v", err)
	}
	return a, nil
}

func parseAssignment(r io.Reader) (*Assignment, error) {
	d := json.NewDecoder(r)
	var asgn assignment
	err := d.Decode(&asgn)
	if err != nil {
		return nil, err
	}
	a := NewAssignment()
	switch {
	case asgn.Code == nil:
		return nil, fmt.Errorf("must have code")
	case Code(*asgn.Code).Validate() != nil:
		c := Code(*asgn.Code)
		return nil, codeErrMsg(c, c.Validate())
	}
	c := Code(*asgn.Code)
	a.code = &c
	if asgn.Name != nil {
		a.name = *asgn.Name
	}
	var problems []problemInterface
	for _, p := range asgn.Problems {
		problems = append(problems, p)
	}
	var handins []handinInterface
	for _, h := range asgn.Handins {
		handins = append(handins, h)
	}
	if err := validateProblemsAndHandins(handins, problems); err != nil {
		return nil, err
	}
	var h []Handin
	for _, hh := range asgn.Handins {
		h = append(h, hh.toHandin())
	}
	a.SetHandinsNoValidate(h)
	var p []Problem
	for _, pp := range asgn.Problems {
		p = append(p, pp.toProblem())
	}
	a.SetProblemsNoValidate(p)
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return a, nil
}

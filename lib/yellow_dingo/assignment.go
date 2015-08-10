package yellow_dingo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/synful/kudos/lib/purple_unicorn"
)

type handin struct {
	Code     *string    `json:"code"`
	Due      *time.Time `json:"due"`
	Problems []string   `json:"problems"`
}

type problem struct {
	Code        *string   `json:"code"`
	Name        *string   `json:"name"`
	Points      *float64  `json:"points"`
	Subproblems []problem `json:"subproblems"`
}

type assignment struct {
	Code     *string   `json:"code"`
	Name     *string   `json:"name"`
	Handins  []handin  `json:"handins"`
	Problems []problem `json:"problems"`
}

func ParseAssignmentFile(path string) (*purple_unicorn.Assignment, error) {
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

func ParseAssignment(r io.Reader) (*purple_unicorn.Assignment, error) {
	a, err := parseAssignment(r)
	if err != nil {
		return nil, fmt.Errorf("parse: %v", err)
	}
	return a, nil
}

func parseAssignment(r io.Reader) (*purple_unicorn.Assignment, error) {
	d := json.NewDecoder(r)
	var asgn assignment
	err := d.Decode(&asgn)
	if err != nil {
		return nil, err
	}
	a := purple_unicorn.NewAssignment()
	if asgn.Code != nil {
		a.SetCodeNoValidate(purple_unicorn.Code(*asgn.Code))
	}
	if asgn.Name != nil {
		a.SetName(*asgn.Name)
	}
	var h []purple_unicorn.Handin
	for _, hh := range asgn.Handins {
		if hh.Due == nil {
			if hh.Code == nil {
				return nil, fmt.Errorf("handin must have due date")
			}
			return nil, fmt.Errorf("handin %v must have due date", *hh.Code)
		}
		var hhh purple_unicorn.Handin
		if hh.Code != nil {
			hhh.Code = purple_unicorn.Code(*hh.Code)
			hhh.HasCode = true
		}
		for _, p := range hh.Problems {
			hhh.Problems = append(hhh.Problems, purple_unicorn.Code(p))
		}
		h = append(h, hhh)
	}
	a.SetHandinsNoValidate(h)

	var convertProblems func(p []problem) ([]purple_unicorn.Problem, error)
	convertProblems = func(p []problem) ([]purple_unicorn.Problem, error) {
		var ret []purple_unicorn.Problem
		for _, pp := range asgn.Problems {
			switch {
			case pp.Code == nil:
				return nil, fmt.Errorf("all problems must have codes")
			case pp.Points == nil:
				return nil, fmt.Errorf("problem %v must have points", *pp.Code)
			}
			var ppp purple_unicorn.Problem
			ppp.Code = purple_unicorn.Code(*pp.Code)
			ppp.Points = *pp.Points
			if pp.Name != nil {
				ppp.Name = *pp.Name
			}
			var err error
			ppp.Subproblems, err = convertProblems(pp.Subproblems)
			if err != nil {
				return nil, err
			}
			ret = append(ret, ppp)
		}
		return ret, nil
	}
	p, err := convertProblems(asgn.Problems)
	if err != nil {
		return nil, err
	}
	a.SetProblemsNoValidate(p)
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return a, nil
}

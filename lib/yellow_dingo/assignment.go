package yellow_dingo

import (
	"encoding/json"
	"fmt"
	"github.com/synful/kudos/lib/purple_unicorn"
	"os"
	"time"
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

func ParseAssignment(conf string) (*purple_unicorn.Assignment, error) {
	f, err := os.Open(conf)
	if err != nil {
		return nil, fmt.Errorf("parse %v: %v", err)
	}
	d := json.NewDecoder(f)
	var asgn assignment
	err = d.Decode(&asgn)
	if err != nil {
		return nil, fmt.Errorf("parse %v: %v", err)
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
				return nil, fmt.Errorf("parse %v: handin must have due date", conf)
			}
			return nil, fmt.Errorf("parse %v: handin %v must have due date", conf, *hh.Code)
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
				return nil, fmt.Errorf("parse %v: all problems must have codes", conf)
			case pp.Points == nil:
				return nil, fmt.Errorf("parse %v: problem %v must have points", conf, *pp.Code)
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
		return nil, fmt.Errorf("parse %v: %v", conf, err)
	}
	a.SetProblemsNoValidate(p)
	if err := a.Validate(); err != nil {
		return nil, fmt.Errorf("parse %v: %v", conf, err)
	}
	return a, nil
}

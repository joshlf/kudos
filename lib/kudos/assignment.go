package kudos

import "time"

type Assignment struct {
	Code string
	Name string

	Handins []Handin

	Problems []Problem
}

type Handin struct {
	Code     string
	Due      time.Time
	Problems []string
}

type Problem struct {
	Code string
	Name string

	// If this problem has subproblems,
	// Points is the sum of the point
	// values of all subproblems.
	Points      float64
	Subproblems []Problem
}

func FindAssignmentByCode(as []*Assignment, code string) (a *Assignment, ok bool) {
	for _, aa := range as {
		if aa.Code == code {
			return aa, true
		}
	}
	return nil, false
}

// FindHandinByCode searches through a.Handins for a handin
// with the given code. It panics if the code is not valid.
func (a *Assignment) FindHandinByCode(code string) (h Handin, ok bool) {
	if ValidateCode(code) != nil {
		panic("lib/kudos: FindHandinByCode: invalid code")
	}

	for _, h := range a.Handins {
		// we don't need to worry about h.Code being the empty
		// string because code cannot be (it would be invalid
		// and we'd have already panicked)
		if h.Code == code {
			return h, true
		}
	}
	return Handin{}, false
}

func (a *Assignment) FindProblemByCode(code string) (p Problem, ok bool) {
	if ValidateCode(code) != nil {
		panic("lib/kudos: FindProblemByCode: invalid code")
	}

	return findProblemByCode(code, a.Problems)
}

func findProblemByCode(code string, problems []Problem) (Problem, bool) {
	for _, p := range problems {
		if p.Code == code {
			return p, true
		}
		pp, ok := findProblemByCode(code, p.Subproblems)
		if ok {
			return pp, true
		}
	}
	return Problem{}, false
}

// FindProblemPathByCode returns the list of parents
// of the given problem (not including the problem
// itself). The second return value is true iff the
// problem was found.
func (a *Assignment) FindProblemPathByCode(code string) (parents []string, ok bool) {
	if ValidateCode(code) != nil {
		panic("lib/kudos: FindProblemByCode: invalid code")
	}

	return findProblemPathByCode(code, a.Problems)
}

func findProblemPathByCode(code string, problems []Problem) ([]string, bool) {
	for _, p := range problems {
		if p.Code == code {
			return nil, true
		}
		path, ok := findProblemPathByCode(code, p.Subproblems)
		if ok {
			return append([]string{p.Code}, path...), true
		}
	}
	return nil, false
}

func (a *Assignment) TotalPoints() float64 {
	total := 0.0
	for _, p := range a.Problems {
		total += p.Points
	}
	return total
}

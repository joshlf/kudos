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

	RubricCommentTemplate string

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

// SetHandinByCode searches through a.Handins for a handin
// with the given code, overwriting it with h once it is
// found. It panics if the code is not valid or if no such
// handin is found.
func (a *Assignment) SetHandinByCode(code string, h Handin) {
	if ValidateCode(code) != nil {
		panic("lib/kudos: SetHandinByCode: invalid code")
	}

	for i, hh := range a.Handins {
		// we don't need to worry about h.Code being the empty
		// string because code cannot be (it would be invalid
		// and we'd have already panicked)
		if hh.Code == code {
			a.Handins[i] = h
		}
	}
	panic("lib/kudos: SetHandinByCode: no such code")
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

func (a *Assignment) TraverseProblemsPreOrder(f func(p Problem)) {
	for _, p := range a.Problems {
		p.TraversePreOrder(f)
	}
}

func (a *Assignment) TraverseProblemsPostOrder(f func(p Problem)) {
	for _, p := range a.Problems {
		p.TraversePostOrder(f)
	}
}

func (p Problem) TraversePreOrder(f func(p Problem)) {
	var walkFn func(p Problem)
	walkFn = func(p Problem) {
		f(p)
		for _, pp := range p.Subproblems {
			walkFn(pp)
		}
	}
	walkFn(p)
}

func (p Problem) TraversePostOrder(f func(p Problem)) {
	var walkFn func(p Problem)
	walkFn = func(p Problem) {
		for _, pp := range p.Subproblems {
			walkFn(pp)
		}
		f(p)
	}
	walkFn(p)
}

func (a *Assignment) TotalPoints() float64 {
	total := 0.0
	for _, p := range a.Problems {
		total += p.Points
	}
	return total
}

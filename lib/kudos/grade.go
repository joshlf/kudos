package kudos

type AssignmentGrade struct {
	// Grades contains the grade for every
	// problem, including those which are not
	// at the top level. The invariant must
	// be maintained that if a problem has
	// a grade, none of its children have grades.
	Grades map[string]ProblemGrade
}

// Total computes the total number of points given by
// the AssignmentGrade a on the Assignment asgn. If
// a is not a complete grade, Total returns false.
// If a is not a grade for asgn, the behavior of Total
// is undefined (and it will likely panic).
func (a *AssignmentGrade) Total(asgn *Assignment) (grade float64, ok bool) {
	total := 0.0
	for _, p := range asgn.Problems {
		g, ok := a.ProblemTotal(asgn, p.Code)
		if !ok {
			return 0.0, false
		}
		total += g
	}
	return total, true
}

// ProblemTotal computes the total number of points
// given by the AssignmentGrade a on the given problem
// of the given assignment. If a is not a complete
// grade for the given problem, ProblemTotal returns
// false. If a is not a grade for asgn, the behavior
// of ProblemTotal is undefined (and it will likely panic).
func (a *AssignmentGrade) ProblemTotal(asgn *Assignment, problem string) (grade float64, ok bool) {
	if g, ok := a.Grades[problem]; ok {
		return g.Grade, true
	}
	total := 0.0
	p, _ := asgn.FindProblemByCode(problem)
	if len(p.Subproblems) == 0 {
		return 0.0, false
	}
	for _, pp := range p.Subproblems {
		g, ok := a.ProblemTotal(asgn, pp.Code)
		if !ok {
			return 0.0, false
		}
		total += g
	}
	return total, true
}

type ProblemGrade struct {
	Grade   float64
	Comment string
}

package kudos

type AssignmentGrade struct {
	// Grades contains the grade for every
	// problem, including those which are not
	// at the top level. The invariant must
	// be maintained that if a problem has
	// a grade, none of its children have grades.
	Grades map[string]ProblemGrade
}

func (a *AssignmentGrade) Total() float64 {
	total := 0.0
	for _, g := range a.Grades {
		total += g.Grade
	}
	return total
}

type ProblemGrade struct {
	Grade   float64
	Comment string
}

func GradeComplete(a *Assignment, g *AssignmentGrade) bool {
	// TODO(joshlf)
	panic("unimplemented")
}

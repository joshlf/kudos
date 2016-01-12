package kudos

type AssignmentGrade struct {
	Problems map[string]ProblemGrade
}

type ProblemGrade struct {
	// If len(Subproblems) > 0, the Grade
	// field should be ignored
	Grade       float64
	Subproblems map[string]ProblemGrade
}

func GradeComplete(a Assignment, g AssignmentGrade) bool {
	return len(GradesMissing(a, g)) == 0
}

// GradeMissing returns a slice of problem codes
// for which grades are missing; grades for these
// problems or parent problems must be added in
// order for g to be a complete grade.
func GradesMissing(a Assignment, g AssignmentGrade) []string {
	panic("unimplemented")
	// TODO(joshlf)
}

package kudos

type AssignmentGrade struct {
	// Grades contains the grade for every
	// problem, including those which are not
	// at the top level. The invariant must
	// be maintained that if a problem has
	// a grade, none of its children have grades.
	Grades map[string]ProblemGrade
}

type ProblemGrade struct {
	Grade   float64
	Comment string
}

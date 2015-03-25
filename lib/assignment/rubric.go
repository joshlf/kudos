package assignment

type Rubric struct {
	Assignment string
	Grader     string
	Grades     []Grade `toml:"grade"`
}

type Grade struct {
	Problem string
	Comment string
	Score   GradeNum
	Total   GradeNum `toml:"possible"`
}

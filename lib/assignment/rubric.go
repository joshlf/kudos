package assignment

type Rubric struct {
	Assignment string
	Grader     string
	Grade      []Grade `toml:"grade"`
}

type Grade struct {
	Problem  string   `toml:"problem"`
	Comment  string   `toml:"comment"`
	Score    GradeNum `toml:"score"`
	Possible GradeNum `toml:"possible"`
}

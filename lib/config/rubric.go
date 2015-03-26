package config

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

func DefaultGrade() Grade {
	return Grade{
		Problem:  "Problem Name",
		Comment:  "Comments about the problem",
		Score:    GradeNum(0),
		Possible: GradeNum(100),
	}
}

func DefaultRubric() Rubric {
	return Rubric{
		Assignment: "Assignment Name",
		Grader:     "Grader Login",
		Grade:      []Grade{DefaultGrade()},
	}
}

package config

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/synful/kudos/lib/perm"
)

// WriteRubricFile is equivalent to WriteRubricFileProblems
// for all problems in a.
func WriteRubricFile(a Assignment, file string) error {
	var problems []string
	for _, p := range a.Problems() {
		problems = append(problems, p.Code)
	}
	return WriteRubricFileProblems(a, problems, file)
}

// WriteRubricFileProblems writes creates a template rubric
// for the given problems in the location specified by file.
// The rubric will have all fields pre-populated except for
// scores. Comments are pre-populated with the empty string.
// Unless scores are added, the rubric file will not be valid
// TOML. Thus, a grader will be prevented from accidentally
// submitting a rubric which has not been fully filled in.
//
// WriteRubricFileProblems will panic if len(problems) == 0.
func WriteRubricFileProblems(a Assignment, problems []string, file string) error {
	if len(problems) == 0 {
		panic("config: no problems specified")
	}

	probMap := mapFromProbs(a.Problems())

	var tmpl rubricTemplate
	tmpl.AssignmentCode = a.Code()
	for _, p := range problems {
		pp, ok := probMap[p]
		if !ok {
			panic("config: unkown problem code")
		}
		tmpl.Grades = append(tmpl.Grades, problemToGradeTemplate(pp))
	}

	// WARNING: HORRIBLE HACK
	// Will be removed once the toml package supports this
	// What this does is look for lines of the form:
	//   score = "x"
	// Where there is optional preceding whitespace, and x
	// is either a number or the word "none". If it is "none",
	// the line is simply removed. Otherwise, it is transformed
	// into:
	//   score = # out of x points
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	err := enc.Encode(tmpl)
	if err != nil {
		panic(fmt.Errorf("config: internal error: unexpected toml encoding error: %v", err))
	}
	s := bufio.NewScanner(&buf)
	var obuf bytes.Buffer
	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		if len(fields) != 3 || fields[0] != "score" {
			_, err := fmt.Fprintln(&obuf, line)
			if err != nil {
				panic(fmt.Errorf("config: internal error: unexpected write error: %v", err))
			}
			continue
		}
		if fields[2] == `"none"` {
			continue
		}
		prefix := line[:strings.Index(line, "score")]
		score := fields[2][1 : len(fields[2])-1]
		_, err := fmt.Fprintf(&obuf, "%vscore = # out of %v points\n", prefix, score)
		if err != nil {
			panic(fmt.Errorf("config: internal error: unexpected write error: %v", err))
		}
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm.Parse("rw-rw-rw-"))
	if err != nil {
		return err
	}
	_, err = obuf.WriteTo(f)
	return err
}

func problemToGradeTemplate(p Problem) gradeTemplate {
	var g gradeTemplate
	g.Problem = p.Code
	if len(p.SubProblems) == 0 {
		g.Score = scoreTemplate{p.Points(), true}
	}
	for _, pp := range p.SubProblems {
		g.Grades = append(g.Grades, problemToGradeTemplate(pp))
	}
	return g
}

type rubricTemplate struct {
	AssignmentCode string          `toml:"assignment"`
	Grades         []gradeTemplate `toml:"grade"`
}

type scoreTemplate struct {
	points float64
	set    bool
}

func (s scoreTemplate) MarshalText() ([]byte, error) {
	if !s.set {
		return []byte("none"), nil
	}
	return []byte(fmt.Sprint(s.points)), nil
}

// Rubric represents a filled-in rubric
// for one or more problems on an
// assignment. It is guaranteed that
// this package will never return a
// Rubric which does not contain at
// least one problem or which does not
// contain a grade for every problem
// it contains.
type Rubric struct {
	r rubricFile
}

// AssignmentCode returns the code of the assignment
// that r is for.
func (r Rubric) AssignmentCode() string { return string(r.r.AssignmentCode.code) }

// Grades returns the grades that r contains.
func (r Rubric) Grades() []Grade {
	var g []Grade
	for _, gg := range r.r.Grades {
		g = append(g, gg.toGrade())
	}
	return g
}

// Grade represents a grade on a problem.
type Grade struct {
	Problem string
	Comment string
	score   float64
	Grades  []Grade
}

// Score returns the score of g.
// If len(p.Grades) > 0, it will be inferred
// from the scores of p.Grades. Otherwise, it
// will be the score assigned directly to this
// grade. It is guaranteed that this package
// will never return a Grade which has neither
// a score nor sub-grades, or which has both.
func (g Grade) Score() float64 {
	if len(g.Grades) == 0 {
		return g.score
	}
	score := float64(0)
	for _, gg := range g.Grades {
		score += gg.Score()
	}
	return score
}

type gradeTemplate struct {
	Problem string `toml:"problem"`
	// TODO(synful): are we sure we want to allow comments
	// on problems that don't get grades?
	Comment string          `toml:"comment"`
	Score   scoreTemplate   `toml:"score"`
	Grades  []gradeTemplate `toml:"grade"`
}

type rubricFile struct {
	// Guaranteed to be set
	AssignmentCode optionalCode `toml:"assignment"`
	Grades         []grade      `toml:"grade"`
}

type grade struct {
	// Guaranteed to be set
	Problem optionalCode `toml:"problem"`
	Comment string       `toml:"comment"`

	// Guaranteed that either Score is set
	// or len(Grades) > 0, but not both and
	// not neither.
	Score  optionalNumber `toml:"score"`
	Grades []grade        `toml:"grade"`
}

func (g grade) toGrade() Grade {
	gg := Grade{
		Problem: string(g.Problem.code),
		Comment: g.Comment,
		score:   float64(g.Score.number),
	}
	for _, ggg := range g.Grades {
		gg.Grades = append(gg.Grades, ggg.toGrade())
	}
	return gg
}

// ReadRubricFile reads the rubric at file for the
// assignment a. It is intended that rubrics will
// have been generated by WriteRubricFile or
// WriteRubricFileProblems and subsequently completed
// by a grader.
func ReadRubricFile(a Assignment, file string) (Rubric, error) {
	var rf rubricFile
	if _, err := toml.DecodeFile(file, &rf); err != nil {
		return Rubric{}, err
	}
	if !rf.AssignmentCode.set {
		return Rubric{}, fmt.Errorf("must have assignment code")
	}
	if string(rf.AssignmentCode.code) != a.Code() {
		return Rubric{}, fmt.Errorf("assignment code in rubric (%v) does not match expected code (%v)", string(rf.AssignmentCode.code), a.Code())
	}
	if len(rf.Grades) == 0 {
		return Rubric{}, fmt.Errorf("rubric has no grades")
	}
	var rubric Rubric
	rubric.r = rf
	probs := mapFromProbs(a.Problems())
	for _, g := range rf.Grades {
		if !g.Problem.set {
			fmt.Errorf("all grades must specify problem code")
		}
		_, ok := probs[string(g.Problem.code)]
		if !ok {
			return Rubric{}, fmt.Errorf("unknown problem: %v", g.Problem)
		}
	}
	// TODO(synful): validate subproblems
	return rubric, nil
}

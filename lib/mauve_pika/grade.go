package mauve_pika

import (
	"fmt"

	"github.com/synful/kudos/lib/purple_unicorn"
)

// This file handles grades being entered into the database. In line with the
// database design, kudos conceptually handles grades from the perspective of
// the entire class, with the entire class's grades being committed to the
// database after each change. The grading module does not cary over the
// structural aspects of a grade present in a course config: no
// categories/handins are stored. Instead we rely on the course config
// information to compute grades, enforcing a single point of truth. We
// support:
//     - Storing a student's grade on an assignment
//     - Storing a team's grade on an assignment
//     - Calculating a student's current grade in a class
//     - Changing any of the assignment grades for either a student or a team

// NOTE there is a subtelty in the case of how we handle changes to a course
// config regarding grades. One example is renaming an assignment. Doing this
// could potentially render some grades lost, as we use an assignment name to
// find grades stored in the database. One solution is to add 'rename'
// functionality to kudos
// TODO make grades their own opaque type. Are we committed on float64?

type CourseGrades struct {
	Name        string
	Assignments map[string]AssignmentGrade
	TeamManifest
}

// TODO this implies the constraint that assignment names cannot overlap across
// categories I am not sure if we are already validating that, but it would not
// be hard to add to the validation infrastructure for categories
type StudentGrade struct {
	Team  string
	Score *Score
	// This is an optional field, we may want a more general infrastructure for
	// tagging grades with metadata. For now, this should suffice
	// We store all grading metadata that we can. This could also be a []byte
	Rubrics []Rubric
}

type AssignmentGrade struct {
	Name   string
	Grades map[purple_unicorn.User]StudentGrade
}

// While this is similar in structure to how points are stored in an assignment
// config, there is a key difference. In an assignment config, we can reason
// about subproblems that do not have a score, where grading is done at the
// unit of the closest ancestor that has a point value. Here, because this is
// purely for the purpose of storing grades, if a Score has points, it cannot
// have children. Extra information regarding sub-sub problems isn't understood
// by the grading subsystem, but it is available for inspection, as we store
// all rubrics within an assignment grade.

// NOTE: a simpler design would just store points at the assignment level.
// This, however, makes it more complicated to handle making entries into the
// database one problem at a time. Actually, that may not be the case as each
// entry of points could add the requisite points to the assignment grade. This
// approach certainly makes printing easier though, and it allows us to warn if
// someone is using re-entering a grade for a particular problem

type Score struct {
	HasPoints bool
	Points    float64
	Children  map[string]*Score
}

func (s *Score) Validate() error {
	if s.HasPoints && len(s.Children) != 0 {
		return fmt.Errorf("[internal] a score with points cannot have any children")
	}
	var errs purple_unicorn.ErrList
	for _, c := range s.Children {
		if err := c.Validate(); err != nil {
			errs.Add(err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

type ScorePath []string

func (db *CourseGrades) ChangeTeamGrade(team string, assignment string, path ScorePath, newScore *Score) error {
	t, ok := db.Teams[team]
	if !ok {
		return fmt.Errorf("team %v not found", team)
	}
	for _, student := range t.Members {
		if err := db.ChangeGrade(student, assignment, path, newScore); err != nil {
			return err
		}
	}
	return nil
}

//TODO: singleton pattern?
func (db *CourseGrades) ChangeGrade(student purple_unicorn.User, assignment string, path ScorePath, newScore *Score) error {
	asg, ok := db.Assignments[assignment]
	if !ok {
		return fmt.Errorf("unable to find assignment %v in grades database", assignment)
	}
	grade, ok := asg.Grades[student]
	if !ok {
		return fmt.Errorf("unable to find grade for student '%v' in assignment '%v'", student, assignment)
	}
	reAssign := true
	curGrade := grade.Score
	for _, sc := range path {
		next, ok := curGrade.Children[sc]
		if !ok {
			curGrade.Children[sc] = &Score{HasPoints: false}
			next = curGrade.Children[sc]
		}
		curGrade = next
		reAssign = false
	}
	if reAssign {
		grade.Score = newScore
	} else {
		*curGrade = *newScore
	}
	return nil
}

//Recursively compute the grade for a given score.
//TODO: potential opportunity for code reuse with the assignments code that
//handles a very similar calculation
func (s Score) Grade() float64 {
	if s.HasPoints {
		return s.Points
	}
	res := 0.0
	for _, c := range s.Children {
		res += c.Grade()
	}
	return res
}

//TODO decide how we compute grades. Do we want an average? Or just based on a
//total? Currently, there is just a multiply by weights approach, then dividing
//by total possible given the same weights
//TODO this currenly computes the
//student's final grade. We may want to modify Total() to only compute grades
//for the assignments/hadins that are due, or have been graded.
func (db *CourseGrades) ComputeGrade(student purple_unicorn.User, cat *purple_unicorn.Category) float64 {
	grade := 0.0
	if len(cat.Children()) == 0 {
		for _, asgn := range cat.Assignments() {
			//NOTE assumes that the entries are in the maps
			grade += db.Assignments[asgn.Name()].Grades[student].Score.Grade()
		}
		return grade
	}

	for _, c := range cat.Children() {
		grade += float64(cat.Weight()) * db.ComputeGrade(student, c)
	}
	if possible := cat.Total(); possible == 0 {
		return 0.0
	} else {
		return grade / possible
	}
}

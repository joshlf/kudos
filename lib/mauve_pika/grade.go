package mauve_pika

import "github.com/synful/kudos/lib/purple_unicorn"

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

type CourseGrades struct {
	Name     string
	Students []StudentGrade
}

// TODO this implies the constraint that assignment names cannot overlap across
// categories I am not sure if we are already validating that, but it would not
// be hard to add to the validation infrastructure for categories
type StudentGrade struct {
	Name   string
	Login  purple_unicorn.User
	Grades map[string]AssignmentGrade
}

type AssignmentGrade struct {
	Name  string
	Score Score
	// This is an optional field, we may want a more general infrastructure for
	// tagging grades with metadata. For now, this should suffice
	Team string
	// We store all grading metadata that we can. This could also be a []byte
	Rubrics []Rubric
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
	Children  map[string]Score
}

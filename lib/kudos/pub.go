package kudos

import "time"

type PubDB struct {
	Assignments map[string]*PubAssignment // keys are assignment codes
}

// NewPubDB creates a new PubDB as it should be
// in a newly-initialized course.
func NewPubDB() *PubDB {
	return &PubDB{
		Assignments: make(map[string]*PubAssignment),
	}
}

type PubAssignment struct {
	Code string
	Name string

	Handins []PubHandin
}

type PubHandin struct {
	Code string
	Due  time.Time
}

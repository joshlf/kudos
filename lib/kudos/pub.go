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
	Code   string
	Due    time.Time
	Active bool
}

// FindHandinByCode searches through p.Handins for a handin
// with the given code. It panics if the code is not valid.
func (p *PubAssignment) FindHandinByCode(code string) (h PubHandin, ok bool) {
	if ValidateCode(code) != nil {
		panic("lib/kudos: FindHandinByCode: invalid code")
	}

	for _, h := range p.Handins {
		// we don't need to worry about h.Code being the empty
		// string because code cannot be (it would be invalid
		// and we'd have already panicked)
		if h.Code == code {
			return h, true
		}
	}
	return PubHandin{}, false
}

func AssignmentToPub(a *Assignment) *PubAssignment {
	p := &PubAssignment{
		Code: a.Code,
		Name: a.Name,
	}
	for _, h := range a.Handins {
		// TODO(joshlf): don't default to true
		p.Handins = append(p.Handins, PubHandin{h.Code, h.Due, true})
	}
	return p
}

package kudos

import "time"

type DB struct {
	Students    map[string]*Student    // keys are UIDs
	Assignments map[string]*Assignment // keys are assignment codes
	// keys are assignment codes; values' keys are handin codes,
	// unless there is only one handin, in which case the only
	// key is the empty string; a given assignment/handin's map
	// will exist and be initialized iff the assignment itself
	// is in the Assignments map
	HandinInitialized map[string]map[string]bool
	// keys are assignment codes; value's keys are student UIDs;
	// a given assignment's map will exist and be initialized
	// iff the assignment itself is in the Assignments map
	Grades map[string]map[string]*AssignmentGrade
	// keys are assignment codes; values' keys are handin codes,
	// unless there is only one handin, in which case the only
	// key is the empty string; a given assignment/handin's map
	// will exist and be initialized iff the assignment itself
	// is in the Assignments map
	StudentHandins map[string]map[string]map[string]time.Time

	Anonymizer Anonymizer
}

// AddStudent adds the student with the given uid
// to the database. It returns true if the student
// was added and false if the student already exists
// in the database.
func (d *DB) AddStudent(uid string) bool {
	_, ok := d.Students[uid]
	if ok {
		return false
	}
	d.Students[uid] = &Student{uid}
	return true
}

// AddAssignment adds the given assignment to the
// database. It returns true if the assignment was
// added and false if the assignment already exists
// in the database.
func (d *DB) AddAssignment(a *Assignment) bool {
	_, ok := d.Assignments[a.Code]
	if ok {
		return false
	}
	d.Assignments[a.Code] = a
	d.HandinInitialized[a.Code] = make(map[string]bool)
	d.Grades[a.Code] = make(map[string]*AssignmentGrade)
	d.StudentHandins[a.Code] = make(map[string]map[string]time.Time)
	for _, h := range a.Handins {
		d.HandinInitialized[a.Code][h.Code] = false
		d.StudentHandins[a.Code][h.Code] = make(map[string]time.Time)
	}
	return true
}

// DeleteAssignment deletes the given assignment
// from the database, including all associated
// grades and handins. It returns true if the
// assignment was deleted and false if no
// assignment with the given code was in the
// database.
func (d *DB) DeleteAssignment(code string) bool {
	_, ok := d.Assignments[code]
	if !ok {
		return false
	}
	delete(d.Assignments, code)
	delete(d.HandinInitialized, code)
	delete(d.Grades, code)
	delete(d.StudentHandins, code)
	return true
}

// NewDB creates a new DB as it should be in
// a newly-initialized course.
func NewDB() *DB {
	return &DB{
		Students:          make(map[string]*Student),
		Assignments:       make(map[string]*Assignment),
		HandinInitialized: make(map[string]map[string]bool),
		Grades:            make(map[string]map[string]*AssignmentGrade),
		StudentHandins:    make(map[string]map[string]map[string]time.Time),
		Anonymizer:        NewAnonymizer(),
	}
}

package kudos

type DB struct {
	Students    map[string]*Student    // keys are UIDs
	Assignments map[string]*Assignment // keys are assignment codes
	// keys are assignment codes; value's keys are student UIDs;
	// a given assignment's map will exist and be initialized
	// iff the assignment itself is in the Assignments map
	Grades map[string]map[string]*AssignmentGrade
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
	d.Grades[a.Code] = make(map[string]*AssignmentGrade)
	return true
}

// DeleteAssignment deletes the given assignment
// from the database, including all associated
// grades. It returns true if the assignment was
// deleted and false if no assignment with the
// given code was in the database.
func (d *DB) DeleteAssignment(code string) bool {
	_, ok := d.Assignments[code]
	if !ok {
		return false
	}
	delete(d.Assignments, code)
	delete(d.Grades, code)
	return true
}

// NewDB creates a new DB as it should be in
// a newly-initialized course
func NewDB() *DB {
	return &DB{
		Students:    make(map[string]*Student),
		Assignments: make(map[string]*Assignment),
		Grades:      make(map[string]map[string]*AssignmentGrade),
	}
}

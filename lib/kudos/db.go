package kudos

type DB struct {
	Students map[string]Student // keys are UIDs
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
	d.Students[uid] = Student{uid}
	return true
}

// NewDB creates a new DB as it should be in
// a newly-initialized course
func NewDB() *DB {
	return &DB{
		Students: make(map[string]Student),
	}
}

package kudos

import (
	"path/filepath"

	"github.com/joshlf/kudos/lib/config"
	"github.com/joshlf/kudos/lib/db"
	"github.com/joshlf/kudos/lib/log"
)

type Context struct {
	GlobalConfig *GlobalConfig
	CourseCode   string // if the empty string, considered to be unset
	Course       *Course
	DB           *DB
	PubDB        *PubDB
	committer    db.Committer
	pubcommitter db.Committer
	*log.Logger
}

// OpenDB opens the database, populating the c.DB
// field.
func (c *Context) OpenDB() error {
	d := new(DB)
	committer, err := db.Open(d, c.CourseDBDir())
	if err != nil {
		return err
	}
	c.DB = d
	c.committer = committer
	return nil
}

// CommitDB closes the database, committing
// any changes, and sets the c.DB field to nil.
func (c *Context) CommitDB() error {
	err := c.committer(c.DB)
	c.DB = nil
	c.committer = nil
	return err
}

// CloseDB closes the database without committing
// any changes, and sets the c.DB field to nil.
func (c *Context) CloseDB() error {
	err := c.committer(nil)
	c.DB = nil
	c.committer = nil
	return err
}

// CleanupDB closes the database without committing
// changes if it has not already been closed.
// CleanupDB is meant to be called in defered functions
// or at program exit sites so that any remaining
// database locks are released in case of an unexpected
// exit.
func (c *Context) CleanupDB() error {
	if c.DB != nil {
		return c.CloseDB()
	}
	return nil
}

// OpenPubDB opens the public database, populating the
// c.PubDB field.
func (c *Context) OpenPubDB() error {
	p := new(PubDB)
	committer, err := db.Open(p, c.CoursePubDBDir())
	if err != nil {
		return err
	}
	c.PubDB = p
	c.pubcommitter = committer
	return nil
}

// CommitPubDB closes the public database, committing
// any changes, and sets the c.PubDB field to nil.
func (c *Context) CommitPubDB() error {
	err := c.pubcommitter(c.PubDB)
	c.PubDB = nil
	c.pubcommitter = nil
	return err
}

// ClosePubDB closes the public database without
// committing any changes, and sets the c.PubDB field
// to nil.
func (c *Context) ClosePubDB() error {
	err := c.pubcommitter(nil)
	c.PubDB = nil
	c.pubcommitter = nil
	return err
}

// CleanupPubDB closes the database without committing
// changes if it has not already been closed.
// CleanupPubDB is meant to be called in defered functions
// or at program exit sites so that any remaining
// database locks are released in case of an unexpected
// exit.
func (c *Context) CleanupPubDB() error {
	if c.PubDB != nil {
		return c.ClosePubDB()
	}
	return nil
}

// ReadPubDB reads the public database into the c.PubDB
// field, but does not acquire a lock on it. Thus, no
// changes can be written back to the database (and it
// is an error to call c.CommitPubDB, c.ClosePubDB, or
// c.CleanupPubDB after calling c.ReadPubDB). This method
// is intended to be called when Kudos is being run by
// a non-TA (e.g., a student) who can only read, but not
// alter, the public database.
func (c *Context) ReadPubDB() error {
	p := new(PubDB)
	err := db.Read(p, c.CoursePubDBDir())
	if err != nil {
		return err
	}
	c.PubDB = p
	return nil
}

func (c *Context) CourseRoot() string {
	return CourseCodeToPath(c.CourseCode, c.GlobalConfig)
}

func (c *Context) CourseKudosDir() string {
	return filepath.Join(c.CourseRoot(), config.KudosDirName)
}

func (c *Context) CourseConfigFile() string {
	return filepath.Join(c.CourseKudosDir(), config.CourseConfigFileName)
}

func (c *Context) CourseHandinDir() string {
	return filepath.Join(c.CourseKudosDir(), config.HandinDirName)
}

func (c *Context) CourseSavedHandinsDir() string {
	return filepath.Join(c.CourseKudosDir(), config.SavedHandinsDirName)
}

func (c *Context) CourseAssignmentDir() string {
	return filepath.Join(c.CourseKudosDir(), config.AssignmentDirName)
}

func (c *Context) CourseHooksDir() string {
	return filepath.Join(c.CourseKudosDir(), config.HooksDirName)
}

func (c *Context) PreHandinHookFile() string {
	return filepath.Join(c.CourseHooksDir(), config.PreHandinHookFileName)
}

func (c *Context) CourseDBDir() string {
	return filepath.Join(c.CourseKudosDir(), config.DBDirName)
}

func (c *Context) CoursePubDBDir() string {
	return filepath.Join(c.CourseKudosDir(), config.PubDBDirName)
}

func (c *Context) AssignmentHandinDir(code string) string {
	return filepath.Join(c.CourseHandinDir(), code)
}

func (c *Context) HandinHandinDir(assignment, handin string) string {
	return filepath.Join(c.CourseHandinDir(), assignment, handin)
}

func (c *Context) AssignmentSavedHandinsDir(code string) string {
	return filepath.Join(c.CourseSavedHandinsDir(), code)
}

func (c *Context) HandinSavedHandinsDir(assignment, handin string) string {
	return filepath.Join(c.CourseSavedHandinsDir(), assignment, handin)
}

func (c *Context) UserAssignmentHandinFile(code, uid string) string {
	return filepath.Join(c.AssignmentHandinDir(code), uid, config.HandinFileName)
}

func (c *Context) UserHandinHandinFile(assignment, handin, uid string) string {
	return filepath.Join(c.HandinHandinDir(assignment, handin), uid, config.HandinFileName)
}

func (c *Context) StudentUIDs() []string {
	var uids []string
	for u := range c.DB.Students {
		uids = append(uids, u)
	}
	return uids
}

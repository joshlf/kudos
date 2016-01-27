package main

import (
	"os/user"
	"path/filepath"

	"github.com/joshlf/kudos/lib/config"
	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
)

type student struct {
	// the string that was passed by the user
	// to reference this student; if the user
	// passed a UID, then this is a cleaned
	// UID (leading 0s stripped and so on);
	// if the user passed a username, it is
	// the username
	str     string
	usr     *user.User
	student *kudos.Student
}

func (s *student) String() string { return s.str }

// implements the sort.Interface interface;
// sorting is based on the lexical ordering
// of the str field.
type sortableStudents []*student

func (s sortableStudents) Len() int           { return len(s) }
func (s sortableStudents) Less(i, j int) bool { return s[i].str < s[j].str }
func (s sortableStudents) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Looks up a student by either username or UID,
// and makes sure that they are a student of the
// course. Assumes that the database has been
// opened. If an error is encountered, it is logged
// to ctx.Error, and the process exits.
//
// It is assumed that the argument is obtained
// from a user-supplied command-line argument,
// and the exit codes used are chosen based on this
// assumption.
func lookupStudent(ctx *kudos.Context, u string) *student {
	if len(u) == 0 {
		ctx.Error.Println("bad username or uid: empty")
		exitUsage()
	}

	numeric := true
	for _, c := range u {
		if !(c >= '0' && c <= '9') {
			numeric = false
		}
	}

	var usr *user.User
	var err error

	if numeric {
		usr, err = user.LookupId(u)
		if err != nil {
			ctx.Error.Printf("could not find user with uid %v: %v\n", u, err)
			dev.Fail()
		}
	} else {
		usr, err = user.Lookup(u)
		if err != nil {
			ctx.Error.Printf("could not find user %v: %v\n", u, err)
			dev.Fail()
		}
	}

	s := student{usr: usr}
	if numeric {
		s.str = usr.Uid
	} else {
		s.str = usr.Username
	}

	ss, ok := ctx.DB.Students[usr.Uid]
	if !ok {
		ctx.Error.Printf("no such student: %v\n", s.str)
		exitLogic()
	}
	s.student = ss

	return &s
}

// Looks up the username of the user with the
// given uid, and returns it. If any error is
// encountered, it is logged at the Warn level,
// and the uid is returned instead.
func lookupUsernameForUID(ctx *kudos.Context, uid string) string {
	u, err := user.LookupId(uid)
	if err != nil {
		ctx.Warn.Printf("could not look up username for user with uid %v: %v\n", uid, err)
		return uid
	}
	return u.Username
}

func getUserConfigPath(ctx *kudos.Context) string {
	u, err := user.Current()
	if err != nil {
		ctx.Error.Printf("could not get current user: %v\n", err)
		dev.Fail()
	}
	return filepath.Join(u.HomeDir, config.UserConfigFileName)
}

// attempts to open the database; if an error is
// encountered, it is logged and the process exits
func openDB(ctx *kudos.Context) {
	err := ctx.OpenDB()
	if err != nil {
		ctx.Error.Printf("could not open database: %v\n", err)
		dev.Fail()
	}
}

// attempts to close the database; if an error is
// encountered, it is logged and the process exits
func closeDB(ctx *kudos.Context) {
	err := ctx.CloseDB()
	if err != nil {
		ctx.Error.Printf("could not close database: %v\n", err)
		dev.Fail()
	}
}

// attempts to commit outstanding changes to the
// database; if an error is encountered, it is
// logged and the process exits
func commitDB(ctx *kudos.Context) {
	err := ctx.CommitDB()
	if err != nil {
		ctx.Error.Printf("could not commit changes to database: %v\n", err)
		dev.Fail()
	}
}

// attempts to clean up the database; if an error is
// encountered, it is logged and the process exits
func cleanupDB(ctx *kudos.Context) {
	err := ctx.CleanupDB()
	if err != nil {
		ctx.Error.Printf("could not close database: %v\n", err)
		dev.Fail()
	}
}

// attempts to open the public database; if an error
// is encountered, it is logged and the process exits
func openPubDB(ctx *kudos.Context) {
	err := ctx.OpenPubDB()
	if err != nil {
		ctx.Error.Printf("could not open public database: %v\n", err)
		dev.Fail()
	}
}

// attempts to close the public database; if an error
// is encountered, it is logged and the process exits
func closePubDB(ctx *kudos.Context) {
	err := ctx.ClosePubDB()
	if err != nil {
		ctx.Error.Printf("could not close public database: %v\n", err)
		dev.Fail()
	}
}

// attempts to commit outstanding changes to the
// public database; if an error is encountered,
// it is logged and the process exits
func commitPubDB(ctx *kudos.Context) {
	err := ctx.CommitPubDB()
	if err != nil {
		ctx.Error.Printf("could not commit changes to public database: %v\n", err)
		dev.Fail()
	}
}

// attempts to clean up the public database; if an error
// is encountered, it is logged and the process exits
func cleanupPubDB(ctx *kudos.Context) {
	err := ctx.CleanupPubDB()
	if err != nil {
		ctx.Error.Printf("could not close public database: %v\n", err)
		dev.Fail()
	}
}

// Validates the assignment code and tries to fetch
// the assignment from the database. If either validation
// or lookup fails, an error is logged and the process
// exits (exitUsage for an invalid code and exitLogic for
// a nonexistant assignment). If logName is true, these
// log messages will include the code itself. This is meant
// to be used if there were multiple assignment codes
// specified by the user, and it would be ambiguous not to
// include the code in the log message.
func getAssignment(ctx *kudos.Context, code string, logName bool) *kudos.Assignment {
	validateAssignmentCode(ctx, code, logName)
	asgn, ok := ctx.DB.Assignments[code]
	if !ok {
		if logName {
			ctx.Error.Printf("no such assignment in database: %v\n", code)
		} else {
			ctx.Error.Println("no such assignment in database")
		}
		exitLogic()
	}
	return asgn
}

// Validates the given code. If it is invalid, an error
// is logged and exitUsage() is called. If logCode is true,
// the log message will include the code itself. This is
// meant to be used if there were multiple assignment codes
// specified by the user, and it would be ambiguous not to
// include the code in the log message.
func validateAssignmentCode(ctx *kudos.Context, code string, logCode bool) {
	if err := kudos.ValidateCode(code); err != nil {
		if logCode {
			ctx.Error.Printf("bad assignment code %q: %v\n", code, err)
		} else {
			ctx.Error.Printf("bad assignment code: %v\n", err)
		}
		exitLogic()
	}
}

// Validates the given code. If it is invalid, an error
// is logged and exitUsage() is called. If logCode is true,
// the log message will include the code itself. This is
// meant to be used if there were multiple problem codes
// specified by the user, and it would be ambiguous not to
// include the code in the log message.
func validateProblemCode(ctx *kudos.Context, code string, logCode bool) {
	if err := kudos.ValidateCode(code); err != nil {
		if logCode {
			ctx.Error.Printf("bad problem code %q: %v\n", code, err)
		} else {
			ctx.Error.Printf("bad problem code: %v\n", err)
		}
		exitUsage()
	}
}

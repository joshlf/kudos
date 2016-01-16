package main

import (
	"os/user"

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
// and the uid is returned instead. It is assumed
// that this function is used for looking up
// usernames of students in the database, and
// the log message says this.
func lookupUsernameForUID(ctx *kudos.Context, uid string) string {
	u, err := user.LookupId(uid)
	if err != nil {
		ctx.Warn.Printf("could not look up username for user with uid %v: %v\n", uid, err)
		return uid
	}
	return u.Username
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

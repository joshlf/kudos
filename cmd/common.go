package main

import (
	"os/user"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
)

// prints a standard error and calls dev.Fail on error
func lookupUserByUsernameOrUID(ctx *kudos.Context, u string) *user.User {
	numeric := true
	// consider an empty username to be non-numeric
	if len(u) == 0 {
		numeric = false
	} else {
		for _, c := range u {
			if !(c >= '0' && c <= '9') {
				numeric = false
			}
		}
	}

	// TODO(joshlf): What should we do about empty
	// usernames? It would be weird to use quotation
	// marks to print usernames in the general case,
	// but it would also be weird to behave differently
	// depending on what the username itself is
	if numeric {
		usr, err := user.LookupId(u)
		if err != nil {
			ctx.Error.Printf("could not find user with uid %v: %v\n", u, err)
			dev.Fail()
		}
		return usr
	} else {
		usr, err := user.Lookup(u)
		if err != nil {
			ctx.Error.Printf("could not find user %v: %v\n", u, err)
			dev.Fail()
		}
		return usr
	}
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

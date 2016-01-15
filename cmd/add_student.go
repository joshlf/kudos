package main

import (
	"os/user"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/spf13/cobra"
)

var cmdAddStudent = &cobra.Command{
	Use:   "add-student [usernames]",
	Short: "Add students to the course",
	// TODO(joshlf): long description
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			dev.Fail()
		}
		ctx := getContext()
		addCourseConfig(ctx)

		// Throughout this function, make sure to loop
		// through usernames or uids in the same order
		// as they were given on the command line so
		// that errors are reported in order

		// maps usernames to uids
		uids := make(map[string]string)
		for _, username := range args {
			u, err := user.Lookup(username)
			if err != nil {
				if _, ok := err.(user.UnknownUserError); ok {
					// TODO(joshlf): add "--strict" flag or similar
					// to make these errors (ie, log as error and
					// abort the whole thing)?
					ctx.Warn.Printf("could not find user %v; skipping\n", username)
				} else {
					ctx.Error.Printf("error looking up user %v: %v", username, err)
					dev.Fail()
				}
			} else {
				uids[username] = u.Uid
			}
		}

		err := ctx.OpenDB()
		if err != nil {
			ctx.Error.Printf("could not open database: %v\n", err)
			dev.Fail()
		}
		defer kudos.CleanupDBAndLogOnError(ctx)

		for _, username := range args {
			uid, ok := uids[username]
			if ok {
				ok = ctx.DB.AddStudent(uid)
				if !ok {
					ctx.Warn.Printf("user %v already in database\n", username)
				}
			}
		}

		err = ctx.CommitDB()
		if err != nil {
			ctx.Error.Printf("could not commit changes to database: %v\n", err)
			dev.Fail()
		}
	}
	cmdAddStudent.Run = f
	addAllGlobalFlagsTo(cmdAddStudent.Flags())
	cmdMain.AddCommand(cmdAddStudent)
}

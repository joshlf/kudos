package main

import (
	"fmt"
	"os/user"
	"sort"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/spf13/cobra"
)

var cmdStudent = &cobra.Command{
	Use:   "student",
	Short: "Manage students in the course",
	// TODO(joshlf): long description
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.Usage()
			exitUsage()
		}
		ctx := getContext()
		addCourseConfig(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		var uids []string
		for _, s := range ctx.DB.Students {
			uids = append(uids, s.UID)
		}
		sort.Strings(uids)
		// TODO(joshlf): sort numerically
		//
		// It's mostly important that the order of any
		// warning messages is constant (ie, warning
		// messages produced by calls to lookupUsernameForUID),
		// but it would be nice if these were in numerical
		// as opposed to alphabetical order

		var unames []string
		for _, u := range uids {
			uname := lookupUsernameForUID(ctx, u)
			unames = append(unames, uname)
		}

		sort.Strings(unames)

		for _, u := range unames {
			fmt.Println(u)
		}

	}
	cmdStudent.Run = f
	addAllGlobalFlagsTo(cmdStudent.Flags())
	cmdMain.AddCommand(cmdStudent)
}

var cmdStudentAdd = &cobra.Command{
	Use:   "add [usernames | uids]",
	Short: "Add students to the course",
	// TODO(joshlf): long description
}

func init() {
	var strictFlag bool
	f := func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			exitUsage()
		}
		ctx := getContext()
		addCourseConfig(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		changed := false
		addUser := func(uid, str string) {
			ok := ctx.DB.AddStudent(uid)
			if ok {
				changed = true
			} else {
				ctx.Warn.Printf("user %v already in database\n", str)
			}
		}

		handleErr := func(err error, u string) {
			_, ok1 := err.(user.UnknownUserIdError)
			_, ok2 := err.(user.UnknownUserError)
			if ok1 || ok2 {
				if strictFlag {
					ctx.Error.Printf("could not find user %v; aborting (no changes saved)\n", u)
					closeDB(ctx)
					exitLogic()
				} else {
					ctx.Warn.Printf("could not find user %v; skipping\n", u)
				}
			} else {
				ctx.Error.Printf("could not find user %v: %v\n", u, err)
				dev.Fail()
			}
		}

		for _, u := range args {
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
					handleErr(err, u)
				} else {
					addUser(usr.Uid, usr.Uid)
				}
			} else {
				usr, err = user.Lookup(u)
				if err != nil {
					handleErr(err, u)
				} else {
					addUser(usr.Uid, usr.Username)
				}
			}
		}

		if changed {
			commitDB(ctx)
		} else {
			closeDB(ctx)
		}
	}
	cmdStudentAdd.Run = f
	addAllGlobalFlagsTo(cmdStudentAdd.Flags())
	cmdStudentAdd.Flags().BoolVarP(&strictFlag, "strict", "", false, "if any user is not found, the entire operation is aborted")
	cmdStudent.AddCommand(cmdStudentAdd)
}

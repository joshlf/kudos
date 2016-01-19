package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/handin"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/spf13/cobra"
)

var cmdIngestHandins = &cobra.Command{
	Use:   "ingest-handins <assignment> [<handin>]",
	Short: "Permanently store handins and record handin times in database",
	// TODO(joshlf): long description
}

func init() {
	var studentFlag string
	var allHandinsFlag bool
	var forceFlag bool
	f := func(cmd *cobra.Command, args []string) {
		switch {
		case len(args) < 1 || len(args) > 2:
			cmd.Usage()
			exitUsage()
		case len(args) == 2 && allHandinsFlag:
			fmt.Fprintf(os.Stderr, "cannot specify handin and use --all-handins")
			exitUsage()
		}
		ctx := getContext()
		addCourseConfig(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		asgn, ok := ctx.DB.Assignments[args[0]]
		if !ok {
			ctx.Error.Println("no such assignment in database")
			exitLogic()
		}

		type handinDir struct {
			handin    kudos.Handin
			handinDir string
			savedDir  string
		}
		var handins []handinDir
		if len(args) == 2 {
			h, ok := asgn.FindHandinByCode(args[1])
			if !ok {
				ctx.Error.Printf("no such handin for assignment %v\n")
				exitLogic()
			}
			handins = []handinDir{{
				h,
				ctx.HandinHandinDir(args[0], args[1]),
				ctx.HandinSavedHandinsDir(args[0], args[1]),
			}}
		} else {
			if len(handins) > 1 && allHandinsFlag {
				ctx.Warn.Println("warning: assignment has one handin; --all-handins has no effect")
			}
			for _, h := range asgn.Handins {
				handins = append(handins, handinDir{
					h,
					ctx.HandinHandinDir(args[0], h.Code),
					ctx.HandinSavedHandinsDir(args[0], h.Code),
				})
			}
		}

		oneStudent := cmd.Flag("student").Changed
		var students []*student
		if oneStudent {
			students = []*student{lookupStudent(ctx, studentFlag)}
		} else {
			for _, s := range ctx.DB.Students {
				ss := &student{
					student: s,
					str:     lookupUsernameForUID(ctx, s.UID),
				}
				students = append(students, ss)
			}
			sort.Sort(sortableStudents(students))
		}

		// these will be executed after the database
		// changes have been successfully committed
		var postCommitFuncs []func()

		changed := false
		exitErr := false
		for _, h := range handins {
			if len(handins) > 1 {
				ctx.Verbose.Printf("ingesting handin %v\n", h.handin.Code)
			}

			err := os.MkdirAll(h.savedDir, 0770)
			if err != nil {
				ctx.Error.Printf("could not create save directory: %v; skipping\n", err)
				exitErr = true
				continue
			}

			for _, s := range students {
				// use !oneStudent instead of len(students) > 1
				// since there could actually be only one student
				// in the class, but the --student flag was not
				// given
				if !oneStudent {
					ctx.Verbose.Printf("\t%v\n", s)
				}

				// logPrefix is of one of the following forms:
				//  handin
				//  handin first
				//  handin for bob
				//  handin first for bob
				logPrefix := fmt.Sprintf("handin %v", h.handin.Code)
				if len(handins) == 1 {
					logPrefix = "handin"
				}
				if !oneStudent {
					logPrefix = fmt.Sprintf("%v for %v", logPrefix, s)
				}

				hcode := h.handin.Code
				if len(asgn.Handins) == 1 {
					hcode = ""
				}
				if _, ok := ctx.DB.Handins[asgn.Code][hcode][s.student.UID]; ok {
					if forceFlag {
						ctx.Warn.Printf("warning: %v already ingested; overwriting\n", logPrefix)
					} else {
						ctx.Warn.Printf("warning: %v already ingested; skipping (use --force to overwrite)\n", logPrefix)
						continue
					}
				}

				ok, err := handin.HandedIn(h.handinDir, s.student.UID)
				if err != nil {
					ctx.Error.Printf("could not save %v: %v; skipping\n", logPrefix, err)
					exitErr = true
					continue
				}
				if !ok {
					ctx.Warn.Printf("warning: no %v\n", logPrefix)
					continue
				}

				t, err := handin.HandinTime(h.handinDir, s.student.UID)
				if err != nil {
					ctx.Error.Printf("could not get handin time for %v: %v; skipping\n", logPrefix, err)
					exitErr = true
					continue
				}

				ctx.DB.Handins[asgn.Code][hcode][s.student.UID] = t
				changed = true

				// make sure that all variables used
				// are in local scope so that they
				// are not overwritten on the next
				// loop iteration (since they need
				// to be closed over in the closure's
				// environment)
				handinDir := h.handinDir
				savedDir := h.savedDir
				uid := s.student.UID
				postCommitFuncs = append(postCommitFuncs, func() {
					err := handin.SaveFaclHandin(handinDir, savedDir, uid)
					if err != nil {
						ctx.Error.Printf("could not save %v: %v\n", logPrefix, err)
						exitErr = true
					}
				})
			}
		}

		if changed {
			commitDB(ctx)
		} else {
			closeDB(ctx)
		}
		ctx.Verbose.Println("handin times successfully committed to database; moving handins to permanent storage")

		for _, f := range postCommitFuncs {
			f()
		}

		if exitErr {
			dev.Fail()
		}
	}
	cmdIngestHandins.Run = f
	addAllGlobalFlagsTo(cmdIngestHandins.Flags())
	cmdIngestHandins.Flags().StringVarP(&studentFlag, "student", "", "", "only ingest this student's handin")
	cmdIngestHandins.Flags().BoolVarP(&allHandinsFlag, "all-handins", "", false, "if the assignment has multiple handins, ingest them all")
	cmdIngestHandins.Flags().BoolVarP(&forceFlag, "force", "", false, "overwrite previously-ingested handins")
	cmdMain.AddCommand(cmdIngestHandins)
}

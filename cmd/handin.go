package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sort"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/handin"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/joshlf/kudos/lib/log"
	"github.com/joshlf/kudos/lib/perm"
	"github.com/spf13/cobra"
)

var cmdHandin = &cobra.Command{
	// TODO(joshlf): figure out how to
	// mark the handin as optional in a
	// way that is consistent with Cobra's
	// output.
	Use:   "handin <assignment> [<handin>]",
	Short: "Hand in an assignment",
}

func init() {
	// f := func(cmd *cobra.Command, args []string) {
	// 	ctx := getContext()
	// 	addCourseConfig(ctx)

	// 	var handinFile string
	// 	switch len(args) {
	// 	case 0:
	// 		ctx.Info.Printf("Usage: %v\n\n", cmd.Use)
	// 		asgns, err := kudos.ParseAllAssignmentFiles(ctx)
	// 		if err != nil {
	// 			ctx.Error.Println("could not read all assignments; aborting")
	// 			dev.Fail()
	// 		}
	// 		ctx.Info.Println("Available handins:")
	// 		for _, a := range asgns {
	// 			if len(a.Handins) == 1 {
	// 				ctx.Info.Printf("  %v\n", a.Code)
	// 			} else {
	// 				// TODO(joshlf): maybe change the output
	// 				// format? This works for now, but we
	// 				// could think of something better.
	// 				ctx.Info.Printf("  %v [", a.Code)
	// 				h := a.Handins
	// 				for _, hh := range h[:len(h)-1] {
	// 					ctx.Info.Printf("%v | ", hh.Code)
	// 				}
	// 				ctx.Info.Printf("%v]\n", h[len(h)-1].Code)
	// 			}
	// 		}
	// 		exitClean()
	// 	case 1:
	// 		asgns, err := kudos.ParseAllAssignmentFiles(ctx)
	// 		if err != nil {
	// 			ctx.Error.Println("could not read all assignments; aborting")
	// 			dev.Fail()
	// 		}
	// 		a, ok := kudos.FindAssignmentByCode(asgns, args[0])
	// 		if !ok {
	// 			ctx.Error.Printf("no such assignment: %v\n", args[0])
	// 			exitLogic()
	// 		}
	// 		if len(a.Handins) > 1 {
	// 			// TODO(joshlf): print more useful message,
	// 			// such as available handins?
	// 			ctx.Error.Printf("assignment has multiple handins; please specify one\n")
	// 			exitUsage()
	// 		}
	// 		u, err := user.Current()
	// 		if err != nil {
	// 			ctx.Error.Printf("could not get current user: %v\n", err)
	// 			dev.Fail()
	// 		}
	// 		handinFile = ctx.UserAssignmentHandinFile(args[0], u.Uid)
	// 	case 2:
	// 		asgns, err := kudos.ParseAllAssignmentFiles(ctx)
	// 		if err != nil {
	// 			ctx.Error.Println("could not read all assignments; aborting")
	// 			dev.Fail()
	// 		}
	// 		a, ok := kudos.FindAssignmentByCode(asgns, args[0])
	// 		if !ok {
	// 			ctx.Error.Printf("no such assignment: %v\n", args[0])
	// 			exitLogic()
	// 		}
	// 		_, ok = a.FindHandinByCode(args[1])
	// 		if !ok {
	// 			ctx.Error.Printf("no such handin: %v\n", args[1])
	// 			exitLogic()
	// 		}
	// 		u, err := user.Current()
	// 		if err != nil {
	// 			ctx.Error.Printf("could not get current user: %v\n", err)
	// 			dev.Fail()
	// 		}
	// 		handinFile = ctx.UserHandinHandinFile(args[0], args[1], u.Uid)
	// 	default:
	// 		cmd.Usage()
	// 		exitUsage()
	// 	}

	// 	/*
	// 		Perform handin
	// 	*/

	// 	hook := ctx.PreHandinHookFile()
	// 	doHook := true
	// 	_, err := os.Stat(hook)
	// 	if err != nil {
	// 		if os.IsNotExist(err) {
	// 			doHook = false
	// 		} else {
	// 			ctx.Error.Printf("could not stat pre-handin hook: %v\n", err)
	// 			dev.Fail()
	// 		}
	// 	}

	// 	// TODO(joshlf): Set environment variables

	// 	if doHook {
	// 		c := exec.Command(hook)
	// 		c.Stderr = os.Stderr
	// 		c.Stdout = os.Stdout
	// 		err = c.Run()
	// 		if err != nil {
	// 			if _, ok := err.(*exec.ExitError); ok {
	// 				ctx.Warn.Println("pre-handin hook exited with error code; aborting")
	// 				dev.Fail()
	// 			}
	// 			ctx.Error.Printf("could not run pre-handin hook: %v\n", err)
	// 			dev.Fail()
	// 		}
	// 	}

	// 	printFiles := ctx.Logger.GetLevel() <= log.Info
	// 	if printFiles {
	// 		ctx.Info.Println("Handing in the following files:")
	// 	}
	// 	err = handin.PerformFaclHandin(handinFile, printFiles)
	// 	if err != nil {
	// 		ctx.Error.Printf("could not hand in: %v\n", err)
	// 		dev.Fail()
	// 	}
	// 	ctx.Info.Println("Handin successful.")
	// }
	f := func(cmd *cobra.Command, args []string) {
		ctx := getContext()
		addCourseConfig(ctx)

		var handinFile string
		switch len(args) {
		case 0:
			// TODO(joshlf): Stat the assignment or handin
			// handin directory as a way of telling whether
			// the handin is actually live yet (and maybe
			// make "does the directory exist?" the definition
			// of a handin being live?)
			ctx.Info.Printf("Usage: %v\n\n", cmd.Use)
			readPubDB(ctx)
			var codes []string
			for acode := range ctx.PubDB.Assignments {
				codes = append(codes, acode)
			}
			sort.Strings(codes)

			ctx.Info.Println("Available handins:")
			for _, acode := range codes {
				asgn := ctx.PubDB.Assignments[acode]
				if len(asgn.Handins) == 1 {
					if ctx.Logger.GetLevel() <= log.Verbose {
						ctx.Info.Printf("  %v (due %v)\n", asgn.Code, asgn.Handins[0].Due)
					} else {
						ctx.Info.Printf("  %v\n", asgn.Code)
					}
				} else {
					// TODO(joshlf): maybe change the output
					// format? This works for now, but we
					// could think of something better.
					if ctx.Logger.GetLevel() <= log.Verbose {
						ctx.Info.Printf("  %v\n", asgn.Code)
						for _, h := range asgn.Handins {
							ctx.Info.Printf("    %v (due %v)\n", h.Code, h.Due)
						}
					} else {
						ctx.Info.Printf("  %v [", asgn.Code)
						h := asgn.Handins
						for _, hh := range h[:len(h)-1] {
							ctx.Info.Printf("%v | ", hh.Code)
						}
						ctx.Info.Printf("%v]\n", h[len(h)-1].Code)
					}
				}
			}
			exitClean()
		case 1:
			readPubDB(ctx)
			asgn, ok := ctx.PubDB.Assignments[args[0]]
			if !ok {
				ctx.Error.Printf("no such assignment: %v\n", args[0])
				exitLogic()
			}
			if len(asgn.Handins) > 1 {
				// TODO(joshlf): print more useful message,
				// such as available handins?
				ctx.Error.Printf("assignment has multiple handins; please specify one\n")
				exitUsage()
			}
			u, err := user.Current()
			if err != nil {
				ctx.Error.Printf("could not get current user: %v\n", err)
				dev.Fail()
			}
			handinFile = ctx.UserAssignmentHandinFile(args[0], u.Uid)
		case 2:
			readPubDB(ctx)
			asgn, ok := ctx.PubDB.Assignments[args[0]]
			if !ok {
				ctx.Error.Printf("no such assignment: %v\n", args[0])
				exitLogic()
			}
			_, ok = asgn.FindHandinByCode(args[1])
			if !ok {
				ctx.Error.Printf("no such handin: %v\n", args[1])
				exitLogic()
			}
			u, err := user.Current()
			if err != nil {
				ctx.Error.Printf("could not get current user: %v\n", err)
				dev.Fail()
			}
			handinFile = ctx.UserHandinHandinFile(args[0], args[1], u.Uid)
		default:
			cmd.Usage()
			exitUsage()
		}

		/*
			Perform handin
		*/

		hook := ctx.PreHandinHookFile()
		doHook := true
		_, err := os.Stat(hook)
		if err != nil {
			if os.IsNotExist(err) {
				doHook = false
			} else {
				ctx.Error.Printf("could not stat pre-handin hook: %v\n", err)
				dev.Fail()
			}
		}

		// TODO(joshlf): Set environment variables

		if doHook {
			c := exec.Command(hook)
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			err = c.Run()
			if err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					ctx.Warn.Println("pre-handin hook exited with error code; aborting")
					dev.Fail()
				}
				ctx.Error.Printf("could not run pre-handin hook: %v\n", err)
				dev.Fail()
			}
		}

		printFiles := ctx.Logger.GetLevel() <= log.Info
		if printFiles {
			ctx.Info.Println("Handing in the following files:")
		}
		err = handin.PerformFaclHandin(handinFile, printFiles)
		if err != nil {
			ctx.Error.Printf("could not hand in: %v\n", err)
			dev.Fail()
		}
		ctx.Info.Println("Handin successful.")
	}
	cmdHandin.Run = f
	addAllGlobalFlagsTo(cmdHandin.Flags())
	cmdMain.AddCommand(cmdHandin)
}

var cmdHandinInit = &cobra.Command{
	Use:   "init <assignment>",
	Short: "Initialize an assignment's handins",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			exitUsage()
		}

		ctx := getContext()
		addCourse(ctx)

		if err := kudos.ValidateCode(args[0]); err != nil {
			ctx.Error.Printf("bad assignment code %q: %v\n", args[0], err)
			exitUsage()
		}
		handins := args[1:]
		bad := false
		for _, h := range handins {
			if err := kudos.ValidateCode(h); err != nil {
				ctx.Error.Printf("bad handin code %q: %v\n", h, err)
				bad = true
			}
		}
		if bad {
			exitUsage()
		}

		openDB(ctx)
		defer cleanupDB(ctx)

		asgn, ok := ctx.DB.Assignments[args[0]]
		if !ok {
			ctx.Error.Printf("no such assignment: %v\n", args[0])
			exitLogic()
		}

		bad = false
		for _, h := range handins {
			if _, ok := asgn.FindHandinByCode(h); !ok {
				ctx.Error.Printf("no such handin: %v\n", h)
				bad = true
			}
		}
		if bad {
			exitLogic()
		}

		var uids []string
		for _, s := range ctx.DB.Students {
			uids = append(uids, s.UID)
		}
		closeDB(ctx)

		// If they specified handins, then we've already
		// validated that they exist, so just use those.
		// If they haven't specified handins, it could either
		// be because there aren't any named handins, or
		// because they want to initialize all of them.
		if len(handins) == 0 && len(asgn.Handins) > 1 {
			ctx.Verbose.Println("no handins specified, but assignment has multiple handins; initializing all of them")
			for _, h := range asgn.Handins {
				handins = append(handins, h.Code)
			}
		}
		err := handinInit(ctx, asgn, uids, handins...)
		if err != nil {
			ctx.Error.Printf("initialization failed: %v\n", err)
			dev.Fail()
		}

		// // If there is a single handin, initialize the handin
		// // directory directly. Otherwise, create the parent
		// // directory and initialize each handin directory
		// // one at a time.
		// if len(asgn.Handins) == 1 {
		// 	dir := ctx.AssignmentHandinDir(asgn.Code)
		// 	err := handin.InitFaclHandin(dir, uids)
		// 	if err != nil {
		// 		ctx.Error.Printf("initialization failed: %v", err)
		// 		dev.Fail()
		// 	}
		// } else {
		// 	// need world r-x so students can cd in
		// 	// and write to their handin files
		// 	mode := perm.Parse("rwxrwxr-x")
		// 	dir := ctx.AssignmentHandinDir(asgn.Code)
		// 	err = os.Mkdir(dir, mode)
		// 	if err != nil {
		// 		ctx.Error.Printf("could not create handin directory: %v\n", err)
		// 		dev.Fail()
		// 	}
		// 	// set permissions explicitly since original permissions
		// 	// might be masked (by umask)
		// 	err = os.Chmod(dir, mode)
		// 	if err != nil {
		// 		ctx.Error.Printf("could not set permissions on handin directory: %v\n", err)
		// 		dev.Fail()
		// 	}
		// 	for _, h := range asgn.Handins {
		// 		dir := ctx.HandinHandinDir(asgn.Code, h.Code)
		// 		err := handin.InitFaclHandin(dir, uids)
		// 		if err != nil {
		// 			ctx.Error.Printf("could not initialize handin %v: %v", h.Code, err)
		// 			dev.Fail()
		// 		}
		// 	}
		// }
	}
	cmdHandinInit.Run = f
	addAllGlobalFlagsTo(cmdHandinInit.Flags())
	cmdHandin.AddCommand(cmdHandinInit)
}

// handinInit initializes the given assignment's handins. If the assignment
// has multiple handins, it creates the top-level assignment directory if
// it does not already exist, but assumes that each handin directory specified
// by the handins argument does not exist yet. If the assignment only has one
// handin, it assumes that the assignment directory does not exist yet.
func handinInit(ctx *kudos.Context, asgn *kudos.Assignment, uids []string, handins ...string) error {
	switch {
	case len(asgn.Handins) > 1 && len(handins) == 0:
		panic("internal: no handins specified in argument to handinInit")
	case len(asgn.Handins) == 1 && len(handins) > 0:
		panic("internal: handins spuriously specified in argument to handinInit")
	}
	var h []kudos.Handin
	if len(asgn.Handins) == 0 {
		h = []kudos.Handin{asgn.Handins[0]}
	} else {
		for _, hcode := range handins {
			hh, ok := asgn.FindHandinByCode(hcode)
			if !ok {
				panic("internal: bad handin code given in argument to handinInit")
			}
			h = append(h, hh)
		}
	}

	if len(asgn.Handins) == 1 {
		dir := ctx.AssignmentHandinDir(asgn.Code)
		err := handin.InitFaclHandin(dir, uids)
		if err != nil {
			return err
		}
	} else {
		dir := ctx.AssignmentHandinDir(asgn.Code)
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				// need world r-x so students can cd in
				// and write to their handin files
				mode := perm.Parse("rwxrwxr-x")
				err := perm.Mkdir(dir, mode)
				if err != nil {
					return fmt.Errorf("create handin directory: %v", err)
				}
			} else {
				return err
			}
		}

		for _, hh := range h {
			dir := ctx.HandinHandinDir(asgn.Code, hh.Code)
			err := handin.InitFaclHandin(dir, uids)
			if err != nil {
				return fmt.Errorf("initialize handin %v: %v", hh.Code, err)
			}
		}
	}
	return nil
}

var cmdHandinIngest = &cobra.Command{
	Use:   "ingest <assignment> [<handin>]",
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

		asgn := getAssignment(ctx, args[0], false)

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
	cmdHandinIngest.Run = f
	addAllGlobalFlagsTo(cmdHandinIngest.Flags())
	cmdHandinIngest.Flags().StringVarP(&studentFlag, "student", "", "", "only ingest this student's handin")
	cmdHandinIngest.Flags().BoolVarP(&allHandinsFlag, "all-handins", "", false, "if the assignment has multiple handins, ingest them all")
	cmdHandinIngest.Flags().BoolVarP(&forceFlag, "force", "", false, "overwrite previously-ingested handins")
	cmdHandin.AddCommand(cmdHandinIngest)
}

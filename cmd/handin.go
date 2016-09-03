package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/handin"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/joshlf/kudos/lib/log"
	"github.com/joshlf/kudos/lib/perm"
	"github.com/spf13/cobra"
)

var cmdHandin = &cobra.Command{
	Use:   "handin <assignment> [<handin>]",
	Short: "Hand in an assignment",
}

func init() {
	var handinLateFlag bool
	f := func(cmd *cobra.Command, args []string) {
		ctx := getContext()
		addCourseConfig(ctx)

		var handinFile string
		switch len(args) {
		case 0:
			ctx.Info.Printf("Usage: %v\n\n", cmd.UseLine())
			readPubDB(ctx)
			var codes []string
			for acode := range ctx.PubDB.Assignments {
				codes = append(codes, acode)
			}
			sort.Strings(codes)

			if outIsTerminal() {
				red := color.New(color.FgRed).SprintFunc()("red")
				ctx.Info.Printf("Available handins (handins in %v are inactive):\n", red)
			} else {
				ctx.Info.Println("Available handins:")
			}

			// Note: color package automatically detects whether
			// output file is a tty and disables colorization
			// as needed.
			redSprint := color.New(color.FgRed).SprintFunc()
			formatCode := func(code string, h kudos.PubHandin) string {
				if !h.Active {
					code = redSprint(code)
				}
				return code
			}
			for _, acode := range codes {
				asgn := ctx.PubDB.Assignments[acode]
				if len(asgn.Handins) == 1 {
					code := formatCode(asgn.Code, asgn.Handins[0])
					if ctx.Logger.GetLevel() <= log.Verbose {
						ctx.Info.Printf("  %v (due %v)\n", code, asgn.Handins[0].Due)
					} else {
						ctx.Info.Printf("  %v\n", code)
					}
				} else {
					// TODO(joshlf): maybe change the output
					// format? This works for now, but we
					// could think of something better.
					if ctx.Logger.GetLevel() <= log.Verbose {
						ctx.Info.Printf("  %v\n", asgn.Code)
						for _, h := range asgn.Handins {
							ctx.Info.Printf("    %v (due %v)\n", formatCode(h.Code, h), h.Due)
						}
					} else {
						ctx.Info.Printf("  %v [", asgn.Code)
						h := asgn.Handins
						for _, hh := range h[:len(h)-1] {
							ctx.Info.Printf("%v | ", formatCode(hh.Code, hh))
						}
						hh := h[len(h)-1]
						ctx.Info.Printf("%v]\n", formatCode(hh.Code, hh))
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
				ctx.Error.Println("assignment has multiple handins; please specify one")
				ctx.Error.Println("(use 'kudos handin' to see available handins)")
				exitUsage()
			}
			if !asgn.Handins[0].Active {
				ctx.Error.Println("this handin is inactive; you cannot hand in until it has been activated")
				exitLogic()
			}
			if time.Now().After(asgn.Handins[0].Due) {
				if handinLateFlag {
					ctx.Warn.Println("warning: handing in late")
				} else {
					ctx.Error.Println("handin due date has passed; use --handin-late to hand in anyway")
					exitLogic()
				}
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
			h, ok := asgn.FindHandinByCode(args[1])
			if !ok {
				ctx.Error.Printf("no such handin: %v\n", args[1])
				exitLogic()
			}
			if !h.Active {
				ctx.Error.Println("this handin is inactive; you cannot hand in until it has been activated")
				exitLogic()
			}
			if time.Now().After(h.Due) {
				if handinLateFlag {
					ctx.Warn.Println("warning: handing in late")
				} else {
					ctx.Error.Println("handin due date has passed; use --handin-late to hand in anyway")
					exitLogic()
				}
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
	cmdHandin.Flags().BoolVarP(&handinLateFlag, "--handin-late", "", false, "handin even if it is past the due date")
	cmdMain.AddCommand(cmdHandin)
}

var cmdHandinInit = &cobra.Command{
	Use:   "init <assignment> [<handin>...]",
	Short: "Initialize an assignment's handins",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			exitUsage()
		}

		acode := args[0]
		hcodes := args[1:]

		ctx := getContext()
		addCourse(ctx)
		checkIsTA(ctx)

		validateAssignmentCodes(ctx, acode)
		validateHandinCodes(ctx, hcodes...)

		openDB(ctx)
		defer cleanupDB(ctx)

		asgn := getAssignment(ctx, acode)
		bad := false
		for _, h := range hcodes {
			if _, ok := asgn.FindHandinByCode(h); !ok {
				ctx.Error.Printf("no such handin: %v\n", h)
				bad = true
			}
		}
		if bad {
			exitLogic()
		}

		// If they specified handins, then we've already
		// validated that they exist, so just use those.
		// If they haven't specified handins, it could either
		// be because there aren't any named handins, or
		// because they want to initialize all of them.
		if len(hcodes) == 0 && len(asgn.Handins) > 1 {
			ctx.Verbose.Println("no handins specified, but assignment has multiple handins; initializing all of them")
			for _, h := range asgn.Handins {
				hcodes = append(hcodes, h.Code)
			}
		}
		// manually run through each of the handins so that we
		// can update the database after each one and do a partial
		// commit if need be (ie, if some but not all handins
		// are successfully initialized)
		//
		// at this point, len(handins) == 0 means that the assignment
		// only has one (unnamed) handin
		if len(hcodes) == 0 {
			err := handinInit(ctx, asgn, ctx.StudentUIDs())
			if err != nil {
				ctx.Error.Printf("initialization failed: %v\n", err)
				dev.Fail()
			}
			ctx.DB.HandinInitialized[acode][""] = true
		} else {
			fail := false
			changed := false
			uids := ctx.StudentUIDs()
			for _, h := range hcodes {
				err := handinInit(ctx, asgn, uids, h)
				if err != nil {
					ctx.Error.Printf("initialization failed for handin %v: %v\n", h, err)
					fail = true
				} else {
					ctx.DB.HandinInitialized[acode][h] = true
					changed = true
				}
			}
			if fail {
				if changed {
					err := ctx.CommitDB()
					if err != nil {
						ctx.Error.Printf("could not commit changes to database: %v\n", err)
						f := color.New(color.FgRed).SprintFunc()
						ctx.Error.Println(f("WARNING: some handins were initialized, but are still marked as uninitialized in the database"))
					}
				}
				dev.Fail()
			}
		}
		err := ctx.CommitDB()
		if err != nil {
			ctx.Error.Printf("could not commit changes to database: %v\n", err)
			f := color.New(color.FgRed).SprintFunc()
			if len(hcodes) < 2 {
				// 1 (they specified a single handin) or 0
				// (the assignment has one unnamed handin)
				ctx.Error.Println(f("WARNING: handin was initialized, but is still marked as uninitialized in the database"))
			} else {
				ctx.Error.Println(f("WARNING: handins were initialized, but are still marked as uninitialized in the database"))
			}
			dev.Fail()
		}
	}
	cmdHandinInit.Run = f
	addAllGlobalFlagsTo(cmdHandinInit.Flags())
	addAllTAFlagsTo(cmdHandinInit.Flags())
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

// TODO(joshlf): Give a --force flag for activate
// and deactivate, which sets permissions anyway
// even if the public database already has the
// assignment marked in the way the user wants it.
// This is in case the two get out of sync and the
// user isn't comfortable setting it manually.

var cmdHandinActivate = &cobra.Command{
	Use:   "activate <assignment> [<handin>]",
	Short: "Allow students to submit handins for the given handin",
	// TODO(joshlf): long description
}

func init() {
	// var allHandinsFlag bool
	var initializeFlag bool
	f := func(cmd *cobra.Command, args []string) {
		switch {
		case len(args) < 1 || len(args) > 2:
			cmd.Usage()
			exitUsage()
			// case len(args) == 2 && allHandinsFlag:
			// fmt.Fprintf(os.Stderr, "cannot specify handin and use --all-handins")
			// exitUsage()
		}
		ctx := getContext()

		acode := args[0]
		// getAssignment performs this validation, but we'd
		// like to encounter formatting errors in the order
		// they appear on the command line
		validateAssignmentCodes(ctx, acode)
		var hcode string
		handinSpecified := len(args) == 2
		if handinSpecified {
			hcode = args[1]
			validateHandinCodes(ctx, hcode)
		}

		addCourseConfig(ctx)
		checkIsTA(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		asgn := getAssignment(ctx, acode)
		// var handin kudos.Handin
		if handinSpecified {
			_, ok := asgn.FindHandinByCode(hcode)
			if !ok {
				ctx.Error.Printf("no such handin for assignment %v: %v\n", acode, hcode)
				exitLogic()
			}
		} else if len(asgn.Handins) > 1 {
			ctx.Error.Println("assignment has multiple handins; please specify one")
			exitLogic()
		}

		openPubDB(ctx)
		defer cleanupPubDB(ctx)

		// TODO(joshlf): verify that the assignment has been published,
		// and that the handin has been initialized

		pubasgn := ctx.PubDB.Assignments[acode]
		pubhandin := pubasgn.Handins[0]
		if handinSpecified {
			pubhandin, _ = pubasgn.FindHandinByCode(hcode)
		}

		if pubhandin.Active {
			ctx.Warn.Println("warning: handin is already active")
			exitClean()
		}

		// if we initialize the handin, then we have to mark and commit
		// that change in the database, in which case we need to handle
		// the cleanup logic at the end differently
		handinInitialized := false
		if !ctx.DB.HandinInitialized[acode][hcode] {
			if !initializeFlag {
				ctx.Error.Println("handin hasn't been initialized; run 'kudos handin init' or use --initialize")
				exitLogic()
			} else {
				err := handinInit(ctx, asgn, ctx.StudentUIDs(), hcode)
				if err != nil {
					ctx.Error.Printf("initialize handin: %v\n", err)
					dev.Fail()
				}
				ctx.DB.HandinInitialized[acode][hcode] = true
				handinInitialized = true
			}
		}

		for i, h := range pubasgn.Handins {
			if h.Code == hcode {
				pubasgn.Handins[i].Active = true
			}
		}

		dir := ctx.AssignmentHandinDir(acode)
		if handinSpecified {
			dir = ctx.HandinHandinDir(acode, hcode)
		}
		fi, err := os.Stat(dir)
		if err != nil {
			ctx.Error.Printf("stat handin directory: %v\n", err)
			dev.Fail()
		}
		// just in case the other permissions were changed
		// by TAs manually, respect those changes
		mode := fi.Mode() | perm.ParseSingle("r-x")
		err = os.Chmod(dir, mode)
		if err != nil {
			ctx.Error.Printf("set permissions on handin directory: %v\n", err)
			dev.Fail()
		}
		err = ctx.CommitPubDB()
		if err != nil {
			ctx.Error.Printf("could not commit changes to public database: %v\n", err)
			f := color.New(color.FgRed).SprintFunc()
			ctx.Error.Println(f("WARNING: permissions on handin directory were changed, but the handin is still marked as inactive in the public database"))
			if handinInitialized {
				// this is OK because the only information
				// we've changed is whether the handin was
				// initialized, and it was initialized successfully
				err := ctx.CommitDB()
				if err != nil {
					ctx.Error.Printf("could not commit changes to database: %v\n", err)
					f := color.New(color.FgRed).SprintFunc()
					ctx.Error.Println(f("WARNING: handin directory was initialized, but the handin is still marked as uninitialized in the database"))
				}
			}
			dev.Fail()
		}
		if handinInitialized {
			err := ctx.CommitDB()
			if err != nil {
				ctx.Error.Printf("could not commit changes to database: %v\n", err)
				f := color.New(color.FgRed).SprintFunc()
				ctx.Error.Println(f("WARNING: handin directory was initialized, but the handin is still marked as uninitialized in the database"))
				dev.Fail()
			}
		} else {
			closeDB(ctx)
		}
	}
	cmdHandinActivate.Run = f
	addAllGlobalFlagsTo(cmdHandinActivate.Flags())
	addAllTAFlagsTo(cmdHandinActivate.Flags())
	// cmdHandinIngest.Flags().BoolVarP(&allHandinsFlag, "all-handins", "", false, "if the assignment has multiple handins, activate them all")
	cmdHandinActivate.Flags().BoolVarP(&initializeFlag, "initialize", "", false, "if the handin hasn't been initialized yet, initialize it before activation")
	cmdHandin.AddCommand(cmdHandinActivate)
}

var cmdHandinDeactivate = &cobra.Command{
	Use:   "deactivate <assignment> [<handin>]",
	Short: "Disallow students from submitting handins for the given handin",
	// TODO(joshlf): long description
}

func init() {
	// var allHandinsFlag bool
	f := func(cmd *cobra.Command, args []string) {
		switch {
		case len(args) < 1 || len(args) > 2:
			cmd.Usage()
			exitUsage()
			// case len(args) == 2 && allHandinsFlag:
			// fmt.Fprintf(os.Stderr, "cannot specify handin and use --all-handins")
			// exitUsage()
		}
		ctx := getContext()

		acode := args[0]
		// getAssignment performs this validation, but we'd
		// like to encounter formatting errors in the order
		// they appear on the command line
		validateAssignmentCodes(ctx, acode)
		var hcode string
		handinSpecified := len(args) == 2
		if handinSpecified {
			hcode = args[1]
			validateHandinCodes(ctx, hcode)
		}

		addCourseConfig(ctx)
		checkIsTA(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		asgn := getAssignment(ctx, acode)
		if handinSpecified {
			_, ok := asgn.FindHandinByCode(hcode)
			if !ok {
				ctx.Error.Printf("no such handin for assignment %v: %v\n", acode, hcode)
				exitLogic()
			}
			if !ctx.DB.HandinInitialized[acode][hcode] {
				ctx.Error.Println("handin is not initialized")
				exitLogic()
			}
		} else {
			if len(asgn.Handins) > 1 {
				ctx.Error.Println("assignment has multiple handins; please specify one")
				exitLogic()
			}
			if !ctx.DB.HandinInitialized[acode][""] {
				ctx.Error.Println("handin is not initialized")
				exitLogic()
			}
		}

		openPubDB(ctx)
		defer cleanupPubDB(ctx)

		// TODO(joshlf): verify that the assignment has been published

		pubasgn := ctx.PubDB.Assignments[acode]
		pubhandin := pubasgn.Handins[0]
		if handinSpecified {
			pubhandin, _ = pubasgn.FindHandinByCode(hcode)
		}

		if !pubhandin.Active {
			ctx.Warn.Println("warning: handin is already inactive")
			exitClean()
		}
		for i, h := range pubasgn.Handins {
			if h.Code == hcode {
				pubasgn.Handins[i].Active = false
			}
		}

		dir := ctx.AssignmentHandinDir(acode)
		if handinSpecified {
			dir = ctx.HandinHandinDir(acode, hcode)
		}
		fi, err := os.Stat(dir)
		if err != nil {
			ctx.Error.Printf("stat handin directory: %v\n", err)
			dev.Fail()
		}
		// just in case the other permissions were changed
		// by TAs manually, respect those changes
		mode := fi.Mode() & ^perm.ParseSingle("rwx")
		err = os.Chmod(dir, mode)
		if err != nil {
			ctx.Error.Printf("set permissions on handin directory: %v\n", err)
			dev.Fail()
		}
		err = ctx.CommitPubDB()
		if err != nil {
			ctx.Error.Printf("could not commit changes to public database: %v\n", err)
			f := color.New(color.FgRed).SprintFunc()
			ctx.Error.Println(f("WARNING: permissions on handin directory were changed, but the handin is still marked as active in the public database"))
			dev.Fail()
		}
		closeDB(ctx)
	}
	cmdHandinDeactivate.Run = f
	addAllGlobalFlagsTo(cmdHandinDeactivate.Flags())
	addAllTAFlagsTo(cmdHandinDeactivate.Flags())
	// cmdHandinIngest.Flags().BoolVarP(&allHandinsFlag, "all-handins", "", false, "if the assignment has multiple handins, deactivate them all")
	cmdHandin.AddCommand(cmdHandinDeactivate)
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
		checkIsTA(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		asgn := getAssignment(ctx, args[0])

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
				if _, ok := ctx.DB.StudentHandins[asgn.Code][hcode][s.student.UID]; ok {
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

				ctx.DB.StudentHandins[asgn.Code][hcode][s.student.UID] = t
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
	addAllTAFlagsTo(cmdHandinIngest.Flags())
	cmdHandinIngest.Flags().StringVarP(&studentFlag, "student", "", "", "only ingest this student's handin")
	cmdHandinIngest.Flags().BoolVarP(&allHandinsFlag, "all-handins", "", false, "if the assignment has multiple handins, ingest them all")
	cmdHandinIngest.Flags().BoolVarP(&forceFlag, "force", "", false, "overwrite previously-ingested handins")
	cmdHandin.AddCommand(cmdHandinIngest)
}

package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/spf13/cobra"
)

var cmdGrade = &cobra.Command{
	Use:   "grade <assignment> <problem> <student> [<grade>]",
	Short: "Assign a grade to a student on a problem",
	// TODO(joshlf): long description
}

func init() {
	var commentFlag string
	var deleteFlag bool
	var forceFlag bool
	f := func(cmd *cobra.Command, args []string) {
		switch {
		case len(args) != 3 && len(args) != 4:
			cmd.Usage()
			exitUsage()
		case len(args) == 3 && !deleteFlag:
			fmt.Fprintln(os.Stderr, "must specify grade or use --delete flag")
			exitUsage()
		case len(args) == 4 && deleteFlag:
			fmt.Fprintln(os.Stderr, "cannot specify grade and use --delete flag")
			exitUsage()
		}

		ctx := getContext()
		acode := args[0]
		if err := kudos.ValidateCode(acode); err != nil {
			ctx.Error.Printf("bad assignment code %q: %v\n", acode, err)
			exitUsage()
		}
		pcode := args[1]
		if err := kudos.ValidateCode(pcode); err != nil {
			ctx.Error.Printf("bad problem code %q: %v\n", pcode, err)
			exitUsage()
		}
		student := args[2]

		// only used if --delete is not specified
		var grade float64
		if !deleteFlag {
			var err error
			grade, err = strconv.ParseFloat(args[3], 64)
			if err != nil {
				ctx.Error.Printf("could not parse grade %q: %v\n", args[3], err)
				exitUsage()
			}
		}
		addCourseConfig(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		asgn := getAssignment(ctx, acode, false)

		prob, ok := asgn.FindProblemByCode(pcode)
		if !ok {
			ctx.Error.Printf("assignment %v has no problem with the code %v\n", acode, pcode)
			exitLogic()
		}
		u := lookupStudent(ctx, student)

		if deleteFlag {
			if _, ok := ctx.DB.Grades[acode][u.usr.Uid]; !ok {
				ctx.Error.Println("grade does not exist")
				exitLogic()
			}
			gradesMap := ctx.DB.Grades[acode][u.usr.Uid].Grades
			if _, ok := gradesMap[pcode]; !ok {
				ctx.Error.Println("grade does not exist")
				exitLogic()
			}
			delete(gradesMap, pcode)
		} else {
			cur, err := user.Current()
			if err != nil {
				ctx.Error.Printf("could not get current user: %v\n", err)
				dev.Fail()
			}

			// create this student's assignment grade
			// if it doesn't already exist
			if _, ok := ctx.DB.Grades[acode][u.usr.Uid]; !ok {
				ctx.DB.Grades[acode][u.usr.Uid] = &kudos.AssignmentGrade{
					make(map[string]kudos.ProblemGrade)}
			}
			gradesMap := ctx.DB.Grades[acode][u.usr.Uid].Grades

			// we don't need to worry about the "ok" return value;
			// we've already made sure the problem exists
			path, _ := asgn.FindProblemPathByCode(pcode)

			// it doesn't matter what order we traverse the path
			// in because (assuming the database is valid), at
			// most one parent can have a grade assigned to it
			// (if more than one did, that would constitute the
			// same error we're looking for here, and we assume
			// that the database is valid)
			for _, elem := range path {
				if _, ok := gradesMap[elem]; ok {
					ctx.Error.Printf("grade already assigned to parent problem %v\n", elem)
					exitLogic()
				}
			}

			if _, ok := gradesMap[pcode]; ok && !forceFlag {
				ctx.Error.Println("grade already assigned; use --force to overwrite")
				exitLogic()
			}

			// define walkFn even if we don't use it (because --force
			// was passed) since we ust it later
			var walkFn func(p kudos.Problem)
			walkFn = func(p kudos.Problem) {
				for _, pp := range p.Subproblems {
					if _, ok := gradesMap[pp.Code]; ok {
						ctx.Error.Printf("grade already assigned to subproblem %v; use --force to overwrite all subproblem grades\n", pp.Code)
						exitLogic()
					}
					walkFn(pp)
				}
			}
			if !forceFlag {
				walkFn(prob)
			}

			if grade > prob.Points {
				ctx.Warn.Printf("warning: grade is higher than the maximum for this problem (%v points)\n", prob.Points)
			}

			if forceFlag {
				ctx.Warn.Println("warning: overwriting any previous grades for this problem or subproblems")
				walkFn = func(p kudos.Problem) {
					delete(gradesMap, p.Code)
					for _, pp := range p.Subproblems {
						walkFn(pp)
					}
				}
				walkFn(prob)
			}

			gradesMap[pcode] = kudos.ProblemGrade{
				Grade: grade,
				// the zero value of commentFlag is the empty
				// string, so we can just blindly use it
				Comment:   commentFlag,
				GraderUID: cur.Uid,
			}
		}

		commitDB(ctx)
	}
	cmdGrade.Run = f
	addAllGlobalFlagsTo(cmdGrade.Flags())
	cmdGrade.Flags().StringVarP(&commentFlag, "comment", "", "", "the comment associated with this grade")
	cmdGrade.Flags().BoolVarP(&deleteFlag, "delete", "", false, "delete a grade")
	cmdGrade.Flags().BoolVarP(&forceFlag, "force", "f", false, "overwrite previous grade or grades of subproblems")
	cmdMain.AddCommand(cmdGrade)
}

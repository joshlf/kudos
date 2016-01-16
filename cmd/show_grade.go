package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/joshlf/kudos/lib/kudos"
	"github.com/spf13/cobra"
)

var cmdShowGrade = &cobra.Command{
	Use:   "show-grade",
	Short: "Show grades",
	// TODO(joshlf): long description
}

func init() {
	var studentFlag string
	var assignmentFlag string
	var showProblemsFlag bool
	var showTotalsFlag bool
	var precisionFlag uint8

	stripRegex := regexp.MustCompile(`\.?0*$`)
	formatFloat := func(f float64) string {
		fmt.Sprintf("%.*f", int(precisionFlag), f)
		a := []byte(fmt.Sprintf("%.*f", int(precisionFlag), f))
		b := []byte("")
		return string(stripRegex.ReplaceAll(a, b))
	}

	f := func(cmd *cobra.Command, args []string) {
		studentFlagSet := cmdShowGrade.Flags().Lookup("student").Changed
		assignmentFlagSet := cmdShowGrade.Flags().Lookup("assignment").Changed
		switch {
		case len(args) > 0:
			cmd.Usage()
			exitUsage()
		case !studentFlagSet && !assignmentFlagSet:
			fmt.Fprintln(os.Stderr, "must specify --student or --assignment")
			exitUsage()
		}

		ctx := getContext()
		addCourseConfig(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		var asgn *kudos.Assignment
		if assignmentFlagSet {
			if err := kudos.ValidateCode(assignmentFlag); err != nil {
				ctx.Error.Printf("bad assignment code %q: %v\n", assignmentFlag, err)
				exitUsage()
			}
			var ok bool
			asgn, ok = ctx.DB.Assignments[assignmentFlag]
			if !ok {
				ctx.Error.Printf("no assignment in database with the code %v\n", assignmentFlag)
				exitLogic()
			}
		}

		var stud *student
		if studentFlagSet {
			stud = lookupStudent(ctx, studentFlag)
		}

		// if --show-problems was passed, the prefix will be
		// prepended to each line of the output of the per-problem
		// grades
		printGrade := func(uid, assignment, prefix string) {
			grade, ok := ctx.DB.Grades[assignment][uid]
			if !ok || len(grade.Grades) == 0 {
				fmt.Println("missing")
				return
			}
			asgn := ctx.DB.Assignments[assignment]
			total, ok := grade.Total(asgn)
			if ok {
				if showTotalsFlag {
					outOf := asgn.TotalPoints()
					// TODO(joshlf): make the precision variable
					// so that fewer digits are used if they're
					// not needed (ie, trailing 0s removed)
					percent := formatFloat(100 * (total / outOf))
					fmt.Printf("%v/%v (%v%%)\n", formatFloat(total), formatFloat(outOf), percent)
				} else {
					fmt.Println(total)
				}
			} else {
				fmt.Println("incomplete")
			}
			if showProblemsFlag {
				var walkFn func(p kudos.Problem, prefix string)
				walkFn = func(p kudos.Problem, prefix string) {
					// TODO(joshlf): Be able to distinguish between
					// a missing grade (ie, a missing grade on a
					// problem with no subproblems) and an incomplete
					// grade (a problem in which not all subproblems
					// have grades) and give different output.

					total, ok := grade.ProblemTotal(asgn, p.Code)
					// whether this total was calculated from
					// subproblems (as opposed to assigned
					// directly)
					calculated := false
					if _, ok := grade.Grades[p.Code]; !ok {
						calculated = true
					}

					var totalStr string
					if ok {
						if showTotalsFlag {
							// TODO(joshlf): make the precision variable
							// so that fewer digits are used if they're
							// not needed (ie, trailing 0s removed)
							percent := formatFloat(100 * (total / p.Points))
							totalStr = fmt.Sprintf("%v/%v (%v%%)", formatFloat(total), formatFloat(p.Points), percent)
						} else {
							totalStr = formatFloat(total)
						}
					}

					var pointsStr string
					switch {
					case !ok:
						pointsStr = "missing"
					case ok && !calculated:
						pointsStr = fmt.Sprint(totalStr)
					case ok && calculated:
						pointsStr = fmt.Sprintf("%v (calculated from subproblems)", totalStr)
					}

					fmt.Printf("%v%v: %v\n", prefix, p.Code, pointsStr)

					newPrefix := prefix + "\t"
					for _, pp := range p.Subproblems {
						walkFn(pp, newPrefix)
					}
				}
				newPrefix := prefix + "\t"
				for _, p := range asgn.Problems {
					walkFn(p, newPrefix)
				}
			}
		}

		switch {
		case studentFlagSet && assignmentFlagSet:
			fmt.Printf("grade for %v on %v: ", stud.str, asgn.Code)
			printGrade(stud.usr.Uid, asgn.Code, "")
		case studentFlagSet:
			var acodes []string
			for code := range ctx.DB.Assignments {
				acodes = append(acodes, code)
			}
			sort.Strings(acodes)
			for _, code := range acodes {
				fmt.Printf("grade for %v on %v: ", stud.str, code)
				printGrade(stud.usr.Uid, code, "")
			}
		case assignmentFlagSet:
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

			var pairs unameUIDPairs
			for _, u := range uids {
				uname := lookupUsernameForUID(ctx, u)
				pairs = append(pairs, unameUIDPair{uname, u})
			}

			// sort by username
			sort.Sort(pairs)

			for _, pair := range pairs {
				fmt.Printf("grade for %v on %v: ", pair.uname, asgn.Code)
				printGrade(pair.uid, asgn.Code, "")
			}
		}

		closeDB(ctx)
	}
	cmdShowGrade.Run = f
	addAllGlobalFlagsTo(cmdShowGrade.Flags())
	cmdShowGrade.Flags().StringVarP(&studentFlag, "student", "", "", "the student to print grades for")
	cmdShowGrade.Flags().StringVarP(&assignmentFlag, "assignment", "", "", "the assignment to print grades for")
	cmdShowGrade.Flags().BoolVarP(&showProblemsFlag, "show-problems", "", false, "show grade for each problem of an assignment")
	cmdShowGrade.Flags().BoolVarP(&showTotalsFlag, "show-totals", "", false, "show total number of points grades are out of")
	cmdShowGrade.Flags().Uint8VarP(&precisionFlag, "precision", "", 2, "the maximum number of digits of precision to use when formatting floating point values")
	cmdMain.AddCommand(cmdShowGrade)
}

// to make sorting unames alphabetically easier
type unameUIDPair struct {
	uname, uid string
}

type unameUIDPairs []unameUIDPair

func (u unameUIDPairs) Len() int           { return len(u) }
func (u unameUIDPairs) Less(i, j int) bool { return u[i].uname < u[j].uname }
func (u unameUIDPairs) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }

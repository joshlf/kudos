package main

import (
	"os"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/spf13/cobra"
)

var cmdRubric = &cobra.Command{
	Use:   "rubric",
	Short: "Manage rubrics",
	// TODO(joshlf): long description
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		panic("unimplemented")
	}
	cmdRubric.Run = f
	addAllGlobalFlagsTo(cmdRubric.Flags())
	cmdMain.AddCommand(cmdRubric)
}

var cmdRubricGenerate = &cobra.Command{
	Use:   "generate <assignment> <student> <problem> [<problem> [...]]",
	Short: "Generate a rubric",
	// TODO(joshlf): long description
}

func init() {
	var outputFlag string
	var anonymousFlag bool
	f := func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			cmd.Usage()
			exitUsage()
		}
		ctx := getContext()
		addCourseConfig(ctx)

		openDB(ctx)
		defer cleanupDB(ctx)

		asgn := getAssignment(ctx, args[0], false)
		s := lookupStudent(ctx, args[1])

		pcodes := args[2:]

		// maps problems to the most recent
		// problem that they are an ancestor
		// of
		seenParents := make(map[string]string)
		seenCodes := make(map[string]bool)
		for _, code := range pcodes {
			validateProblemCode(ctx, code, true)
			if seenCodes[code] {
				ctx.Error.Printf("duplicate problem code: %v\n", code)
				exitUsage()
			}
			seenCodes[code] = true
			if child, ok := seenParents[code]; ok {
				ctx.Error.Printf("problem %v conflicts with previous problem %v "+
					"(%v is a child of %v)\n", code, child, child, code)
				exitLogic()
			}
			path, ok := asgn.FindProblemPathByCode(code)
			if !ok {
				ctx.Error.Printf("no such problem: %v\n", code)
				exitLogic()
			}
			// it doesn't matter what order we traverse the path
			// in because (assuming the database is valid), at
			// most one parent can have a grade assigned to it
			// (if more than one did, that would constitute the
			// same error we're looking for here, and we assume
			// that the database is valid)
			for _, p := range path {
				if seenCodes[p] {
					ctx.Error.Printf("problem %v conflicts with previous problem %v "+
						"(%v is a child of %v)\n", code, p, code, p)
					exitLogic()
				}
				seenParents[p] = code
			}
		}

		out := os.Stdout
		if cmd.Flag("output").Changed {
			var err error
			out, err = os.Create(outputFlag)
			if err != nil {
				ctx.Error.Printf("could not create output file: %v\n", err)
				dev.Fail()
			}
		}

		var token, uid string
		if anonymousFlag {
			var err error
			token, err = ctx.DB.Anonymizer.NewToken(s.student.UID)
			if err != nil {
				ctx.Error.Printf("could not generate anonymous token: %v\n", err)
				dev.Fail()
			}
			commitDB(ctx)
		} else {
			uid = s.student.UID
			closeDB(ctx)
		}

		err := kudos.GenerateRubric(out, asgn, uid, token, pcodes...)
		if err != nil {
			ctx.Error.Printf("could not generate rubric: %v\n", err)
			dev.Fail()
		}
		err = out.Sync()
		if err != nil {
			ctx.Error.Printf("could not sync output file: %v\n", err)
			dev.Fail()
		}
	}
	cmdRubricGenerate.Run = f
	addAllGlobalFlagsTo(cmdRubricGenerate.Flags())
	cmdRubricGenerate.Flags().StringVarP(&outputFlag, "output", "o", "", "write the rubric to this file instead of stdout")
	cmdRubricGenerate.Flags().BoolVarP(&anonymousFlag, "anonymous", "", false, "store an anonymous token instead of a uid in the rubric")
	cmdRubric.AddCommand(cmdRubricGenerate)
}

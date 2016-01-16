package main

import (
	"os"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/handin"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/joshlf/kudos/lib/perm"
	"github.com/spf13/cobra"
)

var cmdInitHandin = &cobra.Command{
	Use:   "handin <assignment>",
	Short: "Initialize an assignment's handins",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			dev.Fail()
		}

		ctx := getContext()
		addCourse(ctx)

		asgn, err := kudos.ParseAssignment(ctx, args[0])
		if err != nil {
			ctx.Error.Printf("could not read assignment config: %v\n", err)
			dev.Fail()
		}

		openDB(ctx)
		defer cleanupDB(ctx)
		var uids []string
		for _, s := range ctx.DB.Students {
			uids = append(uids, s.UID)
		}
		closeDB(ctx)

		// If there is a single handin, initialize the handin
		// directory directly. Otherwise, create the parent
		// directory and initialize each handin directory
		// one at a time.
		if len(asgn.Handins) == 1 {
			dir := ctx.AssignmentHandinDir(asgn.Code)
			err := handin.InitFaclHandin(dir, uids)
			if err != nil {
				ctx.Error.Printf("initialization failed: %v", err)
				dev.Fail()
			}
		} else {
			// need world r-x so students can cd in
			// and write to their handin files
			mode := perm.Parse("rwxrwxr-x")
			dir := ctx.AssignmentHandinDir(asgn.Code)
			err = os.Mkdir(dir, mode)
			if err != nil {
				ctx.Error.Printf("could not create handin directory: %v\n", err)
				dev.Fail()
			}
			// set permissions explicitly since original permissions
			// might be masked (by umask)
			err = os.Chmod(dir, mode)
			if err != nil {
				ctx.Error.Printf("could not set permissions on handin directory: %v\n", err)
				dev.Fail()
			}
			for _, h := range asgn.Handins {
				dir := ctx.HandinHandinDir(asgn.Code, h.Code)
				err := handin.InitFaclHandin(dir, uids)
				if err != nil {
					ctx.Error.Printf("could not initialize handin %v: %v", h.Code, err)
					dev.Fail()
				}
			}
		}
	}
	cmdInitHandin.Run = f
	addAllGlobalFlagsTo(cmdInitHandin.Flags())
	cmdInit.AddCommand(cmdInitHandin)
}

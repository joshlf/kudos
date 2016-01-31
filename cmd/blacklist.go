package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/joshlf/kudos/lib/perm"
	"github.com/spf13/cobra"
)

var cmdBlacklist = &cobra.Command{
	Use:   "blacklist [usernames | uids]",
	Short: "Manage your blacklist",
	Long: "If invoked with no arguments, blacklist prints the current contents of " +
		"the user's blacklist. Otherwise, it adds/deletes users to/from the blacklist.",
	// TODO(joshlf): long description
}

func init() {
	var deleteFlag bool
	var strictFlag bool
	var privateFlag bool
	f := func(cmd *cobra.Command, args []string) {
		switch {
		case len(args) == 0 && deleteFlag:
			fmt.Fprintln(os.Stderr, "--delete specified; must specify at least one user")
			exitUsage()
		}

		ctx := getContext()
		blacklistPath := getUserBlacklistPath(ctx)
		fi, err := os.Stat(blacklistPath)
		createBlacklist := false
		if err != nil {
			if os.IsNotExist(err) {
				createBlacklist = true
			} else {
				ctx.Error.Printf("could not stat user blacklist: %v\n", err)
				dev.Fail()
			}
		} else {
			if !privateFlag && fi.Mode()&perm.ParseSingle("r--") == 0 {
				ctx.Warn.Println("warning: blacklist is not world-readable; kudos may assign TAs on your blacklist to grade you")
			}
		}

		var uids []string
		// if len(args) == 0, then the user tried to read
		// the contents of a non-existent blacklist; just
		// let the error happen
		//
		// if --delete was specified, then it's an error
		// not to have a blacklist to delete from; just let
		// the error happen
		if !createBlacklist || len(args) == 0 || deleteFlag {
			var err error
			uids, err = kudos.ParseBlacklistFile(blacklistPath)
			if err != nil {
				ctx.Error.Printf("could not read user blacklist: %v\n", err)
				dev.Fail()
			}
		}

		if len(args) == 0 {
			for _, uid := range uids {
				fmt.Println(lookupUsernameForUID(ctx, uid))
			}
		} else {
			changed := false
			handleUser := func(uid, str string) {
				if deleteFlag {
					// delete all instances (instead of just the first)
					// in case the blacklist is malformed and has duplicates
					foundEver := false
					for {
						// repeatedly find the first instance in the list
						// until no instances were found
						found := false
						for i, u := range uids {
							if u == uid {
								copy(uids[i:], uids[i+1:])
								uids = uids[:len(uids)-1]
								found = true
								foundEver = true
								changed = true
								break
							}
						}
						if !found {
							break
						}
					}
					if !foundEver {
						ctx.Warn.Printf("user %v not in blacklist\n", str)
					}
				} else {
					for _, u := range uids {
						if u == uid {
							ctx.Warn.Printf("user %v already in blacklist\n", str)
							return
						}
					}
					uids = append(uids, uid)
					changed = true
				}
			}

			handleErr := func(err error, u string) {
				_, ok1 := err.(user.UnknownUserIdError)
				_, ok2 := err.(user.UnknownUserError)
				if ok1 || ok2 {
					if strictFlag {
						ctx.Error.Printf("could not find user %v; aborting (no changes saved)\n", u)
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
						handleUser(usr.Uid, usr.Uid)
					}
				} else {
					usr, err = user.Lookup(u)
					if err != nil {
						handleErr(err, u)
					} else {
						handleUser(usr.Uid, usr.Username)
					}
				}
			}

			if changed {
				if createBlacklist {
					perms := perm.Parse("rw-rw-r--")
					if privateFlag {
						perms = perm.Parse("rw-rw----")
					}
					// use os.O_EXCL just in case (it's not a very
					// likely race condition, but why not)
					f, err := os.OpenFile(blacklistPath, os.O_CREATE|os.O_EXCL, perms)
					if err != nil {
						ctx.Error.Printf("could not create user blacklist file: %v\n", err)
						dev.Fail()
					}
					f.Close()
				}
				err = kudos.WriteBlacklistFile(blacklistPath, uids...)
				if err != nil {
					ctx.Error.Printf("could not write changes: %v\n", err)
					dev.Fail()
				}
			}
		}
	}
	cmdBlacklist.Run = f
	addAllGlobalFlagsTo(cmdBlacklist.Flags())
	cmdBlacklist.Flags().BoolVarP(&deleteFlag, "delete", "", false, "delete users from your blacklist")
	cmdBlacklist.Flags().BoolVarP(&strictFlag, "strict", "", false, "if any user is not found, the entire operation is aborted")
	cmdBlacklist.Flags().BoolVarP(&privateFlag, "private", "", false, "do not make the blacklist world-readable on creation "+
		"(only useful for TAs who are using the blacklist for grading). Warning: this prevents kudos, when run by another user, "+
		"from reading your blacklist, which means that kudos may assign TAs on your blacklist to grade you.")
	cmdMain.AddCommand(cmdBlacklist)
}

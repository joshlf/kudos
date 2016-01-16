package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/spf13/cobra"
)

var cmdBlacklist = &cobra.Command{
	Use:   "blacklist [usernames | uids]",
	Short: "Manage your blacklist",
	// TODO(joshlf): long description
}

func init() {
	var deleteFlag bool
	var strictFlag bool
	f := func(cmd *cobra.Command, args []string) {
		switch {
		case len(args) == 0 && deleteFlag:
			fmt.Fprintln(os.Stderr, "--delete specified; must specify at least one user")
			exitUsage()
		}

		ctx := getContext()
		configPath := getUserConfigPath(ctx)
		config, err := kudos.ParseUserConfigFile(configPath)
		if err != nil {
			ctx.Error.Printf("could not read user config: %v\n", err)
			dev.Fail()
		}

		if len(args) == 0 {
			for _, uid := range config.Blacklist {
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
						for i, u := range config.Blacklist {
							if u == uid {
								copy(config.Blacklist[i:], config.Blacklist[i+1:])
								config.Blacklist = config.Blacklist[:len(config.Blacklist)-1]
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
					for _, u := range config.Blacklist {
						if u == uid {
							ctx.Warn.Printf("user %v already in blacklist\n", str)
							return
						}
					}
					config.Blacklist = append(config.Blacklist, uid)
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
				err = kudos.WriteUserConfigFile(configPath, config)
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
	cmdMain.AddCommand(cmdBlacklist)
}

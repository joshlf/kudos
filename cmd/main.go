package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/synful/kudos/lib/build"
	"github.com/synful/kudos/lib/config"
	"github.com/synful/kudos/lib/dev"
	"github.com/synful/kudos/lib/log"
)

var cmdMain = &cobra.Command{
	Use:   "kudos",
	Short: "kudos is a simple grading system",
	Long: `kudos is a simple grading system made out of love and frustration by m, ezr,
and jliebowf`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var verboseFlag bool
var quietFlag bool
var debugFlag bool
var configFlag string
var courseFlag string

// globalFlags contains flags which multiple
// different subcommands may use, but which
// the main command will not use.
//
// Subcommands should add these to themselves
// with addGlobalFlagsTo or addAllGlobalFlagsTo
var globalFlags = pflag.NewFlagSet("common", pflag.ContinueOnError)

// Make sure initGlobalFlags is called
// before any inits functions are run
var _ = initGlobalFlags()

// Return a dummy value to enable the above line
func initGlobalFlags() struct{} {
	globalFlags.StringVarP(&configFlag, "config", "", config.DefaultGlobalConfigFile, "location of the global config file")
	globalFlags.StringVarP(&courseFlag, "course", "c", "", "course")
	return struct{}{}
}

// Add the named flags from globalFlags
// to the given flag set
func addGlobalFlagsTo(fset *pflag.FlagSet, flags ...string) {
	for _, fname := range flags {
		f := globalFlags.Lookup(fname)
		if f == nil {
			panic(fmt.Errorf("internal error: unkown named flag: %v", fname))
		}
		fset.AddFlag(f)
	}
}

// Add all flags in globalFlags
// to the given flag set
func addAllGlobalFlagsTo(fset *pflag.FlagSet) {
	globalFlags.VisitAll(func(f *pflag.Flag) { fset.AddFlag(f) })
}

func main() {
	// These flags are used by all subcommands,
	// but are also used by the main command
	// itself (so define them directly on
	// cmdMain rather than putting them in
	// globalFlags)
	cmdMain.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "be more verbose than normal; overrides --quiet")
	cmdMain.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "be more quiet than normal")
	cmdMain.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "print internal debugging information; implies --verbose")
	if build.DebugMode {
		cmdMain.DebugFlags()
	}
	cmdMain.Execute()
}

func common() {
	// If we're in debug mode, leave
	// debug logging on
	if !build.DebugMode {
		if debugFlag {
			log.SetLoggingLevel(log.Debug)
		} else if verboseFlag {
			log.SetLoggingLevel(log.Verbose)
		} else if quietFlag {
			log.SetLoggingLevel(log.Warn)
		}
	} else {
		if verboseFlag {
			log.Debug.Println("debug mode enabled; ignoring --verbose flag")
		}
		if quietFlag {
			log.Debug.Println("debug mode enabled; ignoring --quiet flag")
		}
	}

	err := config.InitConfig(globalFlags.Lookup("config"), globalFlags.Lookup("course"))
	if err != nil {
		log.Error.Printf("could not initialize configuration: %v\n", err)
		dev.Fail()
	}
}

func requireCourse() {
	if config.Config.Course == "" {
		log.Error.Printf("no course provided; please speficy one using the --course flag or the %v%v environment variable\n", config.EnvPrefix, config.CourseEnvVar)
		dev.Fail()
	}
}

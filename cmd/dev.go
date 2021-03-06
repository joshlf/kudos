package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/joshlf/kudos/lib/build"
	"github.com/joshlf/kudos/lib/config"
	"github.com/spf13/cobra"
)

var cmdDev = &cobra.Command{
	Use:   "dev",
	Short: "Print information useful to development of kudos",
	// TODO(joshlf): long description
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		fmt.Println("Printing developer information...")
		fmt.Println("[CONSTANTS]")
		printConstants()
	}
	cmdDev.Run = f
	cmdMain.AddCommand(cmdDev)
}

func printConstants() {
	consts := map[string]interface{}{
		"lib/build.DebugMode":                build.DebugMode,
		"lib/build.DevMode":                  build.DevMode,
		"lib/build.Root":                     build.Root,
		"lib/config.DefaultGlobalConfigFile": config.DefaultGlobalConfigFile,
		"lib/config.KudosDirName":            config.KudosDirName,
		"lib/config.CourseConfigFileName":    config.CourseConfigFileName,
		"lib/config.HandinDirName":           config.HandinDirName,
		"lib/config.AssignmentDirName":       config.AssignmentDirName,
	}
	keys := make([]string, 0)
	for k := range consts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	longest := 0
	for _, k := range keys {
		// Don't push everything off the screen
		// because of a few abnormally long keys
		if len(k) > longest && len(k) < 70 {
			longest = len(k)
		}
	}
	for _, k := range keys {
		diff := longest - len(k)
		spaces := strings.Repeat(" ", diff) + " "
		fmt.Printf("  %v:"+spaces+"%v\n", k, consts[k])
	}

}

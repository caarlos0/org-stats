package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints org-stats version",
	Run: func(cmd *cobra.Command, args []string) {
		if info, ok := debug.ReadBuildInfo(); ok {
			sum := info.Main.Sum
			if sum == "" {
				sum = "none"
			}
			fmt.Printf("https://%s %s @ %s\n", info.Main.Path, info.Main.Version, sum)
		} else {
			fmt.Println("unknown")
		}
	},
}

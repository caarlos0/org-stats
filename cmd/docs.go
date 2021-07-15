package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:                   "docs",
	Short:                 "Generates donuts's command line docs",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Hidden:                true,
	Args:                  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Root().DisableAutoGenTag = true
		return doc.GenMarkdownTree(cmd.Root(), "docs")
	},
}

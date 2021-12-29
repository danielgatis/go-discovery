package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "discovery",
	SilenceUsage: true,
}

// AddCmd adds a command to the root command.
func AddCmd(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

// Execute executes the root command.
func Execute() {
	rootCmd.Execute()
}

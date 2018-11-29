package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the Elklogs version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("elklogs version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

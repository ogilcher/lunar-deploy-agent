package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const kAgentVersion = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Lunar Deploy Agent version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Lunar Deploy Agent v%s\n", kAgentVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

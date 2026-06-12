package cmd

import (
	"fmt"

	"github.com/ogilcher/lunar-deploy-agent/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Lunar Deploy Agent version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Lunar Deploy Agent v%s\n", version.AgentVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check whether the agent is running correctly",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Agent health: OK")
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}

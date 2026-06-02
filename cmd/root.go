package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lunar-client",
	Short: "Lunar Deploy Agent",
	Long:  "A lightweight deployment orchestration agent.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use --help to view available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version displays the version and the build date.",
	Long:  `version displays the version and the build date.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shaloc %s (%s)\n", Version, BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shaloc",
	Short: "SHAre files LOCally !",
	Long: `SHAre files LOCally !
	
shaloc is a tool designed to share files on a local network over HTTP in command line.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

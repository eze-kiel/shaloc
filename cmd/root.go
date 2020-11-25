package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

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
	// Implement graceful shutdown
	var gracefulStop = make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// SIG detection goroutine
	go func() {
		select {
		case sig := <-gracefulStop:
			logrus.Infof("%s received. Exiting...\n", sig)
			os.Exit(0)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

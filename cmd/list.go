package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list lists all the versions that are available",
	Long: `list lists all the versions that are available. For example:

  shaloc update list`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		r, err := parseReleasesAPI()
		if err != nil {
			logrus.Fatal(err)
		}
		fmt.Println("Available versions:")
		displayAvailableVersions(r)
	},
}

func init() {
	updateCmd.AddCommand(listCmd)
}

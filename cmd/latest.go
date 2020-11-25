package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// latestCmd represents the latest command
var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: "latest updates shaloc to the latest version",
	Long: `latest updates shaloc to the latest version. For example:

  shaloc update latest`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		r, err := parseReleasesAPI()
		if err != nil {
			logrus.Fatal(err)
		}
		if err := getLatest(r); err != nil {
			logrus.Errorf("%s", err)
		}
	},
}

func init() {
	updateCmd.AddCommand(latestCmd)
}

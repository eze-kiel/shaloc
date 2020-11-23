package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all the zip files from /tmp/",
	Long: `clean removes all the zip files in /tmp/, even those that were not created by shaloc. For example:

This will remove all the .zip files in /tmp/
  shaloc clean
`,
	Run: func(cmd *cobra.Command, args []string) {
		dirname := "/tmp/"

		// Open the directory and read all its files.
		dirRead, _ := os.Open(dirname)
		dirFiles, _ := dirRead.Readdir(0)

		// Loop over files
		for index := range dirFiles {
			fileHere := dirFiles[index]

			// Get name of file and its full path.
			nameHere := fileHere.Name()
			fullPath := dirname + nameHere

			if string(nameHere[len(nameHere)-4:]) == ".zip" {
				os.Remove(fullPath)
				logrus.Warnf("Wiped %s", fullPath)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download a file from an URL",
	Long: `You can download a file from an URL. For example:

This will create a file called 'out':
  shaloc get -u http://192.168.1.133/file.txt

This will create a file called new.txt:
  shaloc get -u http://192.168.1.133/file.txt -n new.txt
`,

	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		name, _ := cmd.Flags().GetString("name")

		if url == "" {
			fmt.Println("You must provide a URL with the flag -u !")
			os.Exit(1)
		}

		err := download(name, url)
		if err != nil {
			panic(err)
		}

		fmt.Println("Downloaded: " + url)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("url", "u", "", "URL to download the file from.")
	getCmd.Flags().StringP("name", "n", "out", "Name of the file that will be downloaded.")
}

func download(filepath string, url string) error {

	// Init and start the spinner
	s := spinner.New(spinner.CharSets[23], 100*time.Millisecond)
	s.Start()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	// Stop the spinner
	s.Stop()

	return err
}

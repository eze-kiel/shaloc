package cmd

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt is useful if you forget the --aes flag while getting a file.",
	Long: `decrypt allow you to decrypt a file on the disk if you forgot to use --aes flag. For example:

This will decrypt toto.txt:
  shaloc decrypt toto.txt
	
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Type decryption key:\n")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			logrus.Fatalf("%s", err)
		}

		// Init and start the spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Start()
		decryptFile(string(bytePassword), args[0])
		os.Remove(args[0])
		s.Stop()

		fmt.Printf("Decrypted %s in %s.dec\n", args[0], args[0])
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}

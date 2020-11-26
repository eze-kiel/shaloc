package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
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
		output, _ := cmd.Flags().GetString("output")
		useAES, _ := cmd.Flags().GetBool("aes")

		if url == "" {
			fmt.Println("You must provide a URL with the flag -u !")
			os.Exit(1)
		}

		// If no output name is provided, take the last part of the URI
		if output == "" {
			parts := strings.SplitAfter(url, "/")
			output = parts[len(parts)-1]
		}

		// Ask for the passphrase if needed
		var bytePassword []byte
		var err error
		if useAES {
			bytePassword, err = askForPass()
			if err != nil {
				logrus.Fatalf("%s", err)
			}
		}
		err = download(output, url)
		if err != nil {
			logrus.Errorf("%s\n", err)
			return
		}

		fmt.Println("Downloaded: " + output + " from " + url)

		// If --aes, decrypt the file
		if useAES {
			// Init and start the spinner
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Start()

			tmp, err := decryptFile(string(bytePassword), output)
			if err != nil {
				logrus.Fatalf("%s", err)
			}

			if err := os.Remove(output); err != nil {
				logrus.Fatalf("%s", err)
			}

			if err := os.Rename(tmp, output); err != nil {
				logrus.Fatalf("%s", err)
			}

			s.Stop()

			fmt.Printf("Decrypted %s\n.", output)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("url", "u", "", "URL to download the file from.")
	getCmd.Flags().StringP("output", "o", "", "Name of the file that will be downloaded.")
	getCmd.Flags().Bool("aes", false, "Use AES-256 decryption.")
}

// download downloads a file from url and write it in filepath
func download(filepath string, url string) error {

	// Init and start the spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	// Get the data from the url
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

// decryptFile decrypts filename with the key p using AES-256
func decryptFile(p, filename string) (string, error) {
	outFilename := filename + ".dec"

	key := sha256.Sum256([]byte(p))
	ciphertext, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	of, err := os.Create(outFilename)
	if err != nil {
		return "", err
	}
	defer of.Close()

	// cipertext has the original plaintext size in the first 8 bytes, then IV
	// in the next 16 bytes, then the actual ciphertext in the rest of the buffer.
	// Read the original plaintext size, and the IV.
	var origSize uint64
	buf := bytes.NewReader(ciphertext)
	if err = binary.Read(buf, binary.LittleEndian, &origSize); err != nil {
		return "", err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = buf.Read(iv); err != nil {
		return "", err
	}

	// The remaining ciphertext has size=paddedSize.
	paddedSize := len(ciphertext) - 8 - aes.BlockSize
	if paddedSize%aes.BlockSize != 0 {
		return "", fmt.Errorf("want padded plaintext size to be aligned to block size")
	}
	plaintext := make([]byte, paddedSize)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext[8+aes.BlockSize:])

	if _, err := of.Write(plaintext[:origSize]); err != nil {
		return "", err
	}
	return outFilename, nil
}

// askForPass ask for a passphrase twice.
// If they do not match, it returns an error
func askForPass() ([]byte, error) {
	fmt.Println("Type decryption key:")
	try, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}

	fmt.Println("Type decryption key again:")
	try2, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}

	if string(try) != string(try2) {
		return nil, fmt.Errorf("passwords do not match")
	}

	return try, nil
}

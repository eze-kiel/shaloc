package cmd

import (
	"archive/zip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share a file or a folder",
	Long: `share allow you to start a HTTP server to share a file or a folder. For example:

This will share the file test.txt on all interfaces on port 8080
  shaloc share -f test.txt

This will share blah.txt on 192.168.1.36:1337:
  shaloc share -f blah.txt -i 192.168.1.36 -p 1337

This will share the folder /home/user/sup3r-f0ld3r on 0.0.0.0:8080:
  shaloc share -F /home/user/sup3r-f0ld3r
`,

	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		port, _ := cmd.Flags().GetString("port")
		file, _ := cmd.Flags().GetString("file")
		folder, _ := cmd.Flags().GetString("folder")
		randomize, _ := cmd.Flags().GetInt("random")
		maxDownloads, _ := cmd.Flags().GetInt("max")
		useAES, _ := cmd.Flags().GetBool("aes")

		var uri string

		if file == "" && folder == "" {
			fmt.Println("You must provide at least a file to share (-f) or a folder (-F) !")
			os.Exit(1)
		} else if file != "" && folder != "" {
			fmt.Println("You cannot provide a file and a folder !")
			os.Exit(1)
		}

		// If the folder flag is provided...
		if folder != "" {
			// ...be sure to serve a folder
			isFol, err := isFolder(folder)
			if err != nil {
				logrus.Errorf("%s", err)
				return
			}

			// Zip it
			if isFol {
				// If the user provided a full path, we want to keep only the filename.
				parts := strings.Split(folder, "/")
				if parts[len(parts)-1] == "" {
					// If the path ends with /, the last item of parts will be an
					// empty string
					uri = parts[len(parts)-2] + ".zip"
				} else {
					uri = parts[len(parts)-1] + ".zip"
				}

				file, err = compressFolder(folder)
				if err != nil {
					logrus.Errorf("%s", err)
					return
				}
			}
		} else {
			// Check if the file provided is really a file
			isFol, err := isFolder(file)
			if err != nil {
				logrus.Errorf("%s", err)
				return
			}
			// If not, log an error and exit
			if isFol {
				logrus.Errorf("%s", fmt.Errorf("%s is not a file", file))
				os.Exit(1)
			}
			// If the user provided a full path, we want to keep only the filename.
			parts := strings.Split(file, "/")
			uri = parts[len(parts)-1]
		}

		// If the flag -r is provided, randomize the URI
		if randomize > 0 {
			rand.Seed(time.Now().UnixNano())
			uri = randID(randomize)
		}

		// If the flag --aes is provided, ask for a passphrase
		if useAES {
			fmt.Print("Type encryption key:\n")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				logrus.Fatalf("%s", err)
			}

			file, err = encryptFile(string(bytePassword), file)
			if err != nil {
				log.Fatalf("%s", err)
			}
		}

		srv := &http.Server{
			Addr: ip + ":" + port,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		http.HandleFunc("/"+uri, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Disposition", "attachment; filename="+file)
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

			openfile, err := os.Open(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			_, err = io.Copy(w, openfile)
			if err != nil {
				logrus.Errorf("%s", err)

			}

			if maxDownloads >= 0 {
				maxDownloads--
				if maxDownloads == 0 {
					if err := os.Remove(file); err != nil {
						logrus.Fatalf("%s", err)
					}
					cancel()
				}
				logrus.Infof("Downloads remaining: %d", maxDownloads)
			}
		})

		fmt.Printf("Sharing %s on http://%s:%s/%s\n", file, ip, port, uri)
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Warnf("%s", err)
			}
		}()

		select {
		case <-ctx.Done():
			// Shutdown the server when the context is canceled
			if err := srv.Shutdown(ctx); err != nil {
				logrus.Errorf("%s", err)
			}
		}
		logrus.Infof("Max number of downloads reached, shutting down the server.")

	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().StringP("ip", "i", "0.0.0.0", "IP address to serve on.")
	shareCmd.Flags().StringP("port", "p", "8080", "Port to serve on.")
	shareCmd.Flags().StringP("file", "f", "", "File to share.")
	shareCmd.Flags().StringP("folder", "F", "", "Folder to share. It will be zipped.")
	shareCmd.Flags().IntP("random", "r", 0, "Randomize the URI. The integer provided is the random string lentgh.")
	shareCmd.Flags().IntP("max", "m", -1, "Maximum number of downloads.")
	shareCmd.Flags().Bool("aes", false, "Encrypt file with AES-256.")
}

// ifFolder returns true if name is a folder, false elsewhere.
func isFolder(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return false, err
	}

	if mode := fi.Mode(); mode.IsDir() {
		return true, nil
	}

	return false, nil
}

// Generate a random string ID
func randID(n int) string {
	var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

// compressFolder compresses recursively source, and returns the path of the compressed file
func compressFolder(source string) (string, error) {
	// Create a temporary file prefixed with shaloc
	of, err := ioutil.TempFile("", "shaloc")
	if err != nil {
		log.Fatal(err)
	}
	defer of.Close()

	logrus.Infof("Zipping %s into %s...", source, of.Name())

	// Init and start the spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	archive := zip.NewWriter(of)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return "", nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	if err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	}); err != nil {
		log.Fatalf("%s", err)
	}

	s.Stop()

	return of.Name(), err
}

func encryptFile(p, filename string) (string, error) {
	// Create a temporary file prefixed with shaloc
	of, err := ioutil.TempFile("", "shaloc")
	if err != nil {
		log.Fatal(err)
	}
	defer of.Close()

	key := sha256.Sum256([]byte(p))

	plaintext, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Write the original plaintext size into the output file first, encoded in
	// a 8-byte integer.
	origSize := uint64(len(plaintext))
	if err = binary.Write(of, binary.LittleEndian, origSize); err != nil {
		return "", err
	}

	// Pad plaintext to a multiple of BlockSize with random padding.
	if len(plaintext)%aes.BlockSize != 0 {
		bytesToPad := aes.BlockSize - (len(plaintext) % aes.BlockSize)
		padding := make([]byte, bytesToPad)
		if _, err := rand.Read(padding); err != nil {
			return "", err
		}
		plaintext = append(plaintext, padding...)
	}

	// Generate random IV and write it to the output file.
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	if _, err = of.Write(iv); err != nil {
		return "", err
	}

	// Ciphertext has the same size as the padded plaintext.
	ciphertext := make([]byte, len(plaintext))

	// Use AES implementation of the cipher.Block interface to encrypt the whole
	// file in CBC mode.
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	if _, err = of.Write(ciphertext); err != nil {
		return "", err
	}
	return of.Name(), nil
}

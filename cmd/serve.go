package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve a file or a folder",
	Long: `serve allow you to start a HTTP server to serve a file or a folder. For example:

This will serve the file test.txt on 127.0.0.1:8080
  shaloc serve -f test.txt

This will serve blah.txt on 192.168.1.36:1337:
  shaloc serve -f blah.txt -i 192.168.1.36 -p 1337

This will serve the folder /home/user/sup3r-f0ld3r on 127.0.0.1:8080:
  shaloc server -F /home/user/sup3r-f0ld3r
`,

	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		port, _ := cmd.Flags().GetString("port")
		file, _ := cmd.Flags().GetString("file")
		folder, _ := cmd.Flags().GetString("folder")
		randomize, _ := cmd.Flags().GetInt("random")

		var uri string

		if file == "" && folder == "" {
			fmt.Println("You must provide at least a file to share (-f) or a folder (-F) !")
			os.Exit(1)
		}

		if file != "" && folder != "" {
			fmt.Println("You cannot provide a file and a folder !")
			os.Exit(1)
		}

		if folder != "" {
			isFol, err := isFolder(folder)
			if err != nil {
				logrus.Errorf("%s", err)
				return
			}

			// If the provided file is a folder, zip it and share it
			if isFol {
				file, err = compressFolder(folder)
				if err != nil {
					logrus.Errorf("%s", err)
					return
				}
			}
		}

		if randomize > 0 {
			rand.Seed(time.Now().UnixNano())
			uri = randID(randomize)
		} else {
			// If the user provided a full path, we want to keep only the filename.
			parts := strings.Split(file, "/")
			uri = parts[len(parts)-1]
		}

		http.HandleFunc("/"+uri, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Disposition", "attachment; filename="+file)
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

			openfile, err := os.Open(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			io.Copy(w, openfile)
		})

		srv := &http.Server{
			Addr: ip + ":" + port,
		}

		fmt.Printf("Serving %s on http://%s:%s/%s\n", file, ip, port, uri)

		if err := srv.ListenAndServe(); err != nil {
			logrus.Errorf("%s", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("ip", "i", "127.0.0.1", "IP address to serve on.")
	serveCmd.Flags().StringP("port", "p", "8080", "Port to serve on.")
	serveCmd.Flags().StringP("file", "f", "", "File to share.")
	serveCmd.Flags().StringP("folder", "F", "", "Folder to share. It will be zipped.")
	serveCmd.Flags().IntP("random", "r", 0, "Randomize the URI. The integer provided is the random string lentgh.")
}

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

func compressFolder(source string) (string, error) {
	targetFile := "/tmp/" + filepath.Base(source) + ".zip"

	logrus.Infof("Zipping %s into %s...", source, targetFile)

	// Init and start the spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	zipfile, err := os.Create(targetFile)
	if err != nil {
		return "", err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return "", nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
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
	})

	s.Stop()

	return targetFile, err
}

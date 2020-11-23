package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve a file",
	Long: `serve allow you to start a HTTP server to serve a file. For example:

This will serve the file test.txt on 127.0.0.1:8080
  shaloc serve -f test.txt

This will serve blah.txt on 192.168.1.36:1337:
  shaloc serve -f blah.txt -i 192.168.1.36 -p 1337`,
	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		port, _ := cmd.Flags().GetString("port")
		file, _ := cmd.Flags().GetString("file")

		if file == "" {
			fmt.Println("You must provide a file to share with the flag -f !")
			os.Exit(1)
		}

		http.HandleFunc("/"+file, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Disposition", "attachment; filename="+file)
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

			Openfile, err := os.Open(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			io.Copy(w, Openfile)
		})

		fmt.Printf("Now serving on http://%s:%s/%s\n", ip, port, file)
		srv := &http.Server{Addr: ip + ":" + port}

		// Handles networking errors, such as being unable to bind IP or port
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("HTTP Server: ListenAndServe() error: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("ip", "i", "127.0.0.1", "IP address to serve on.")
	serveCmd.Flags().StringP("port", "p", "8080", "Port to serve on.")
	serveCmd.Flags().StringP("file", "f", "", "File to share.")
}

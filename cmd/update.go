package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type releases []struct {
	URL             string    `json:"url"`
	HTMLURL         string    `json:"html_url"`
	ID              int       `json:"id"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		URL                string    `json:"url"`
		ID                 int       `json:"id"`
		NodeID             string    `json:"node_id"`
		Name               string    `json:"name"`
		Label              string    `json:"label"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
}

type system struct {
	os   string
	arch string
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update allow user to update shaloc.",
	Long: `With update, you can update shaloc. For example:

This will update shaloc to the latest release:
  shaloc update latest

This will list all available shaloc versions:
  shaloc update list

This will update shaloc to v1.2.0:
  shaloc update v1.2.0  `,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		r, err := parseReleasesAPI()
		if err != nil {
			logrus.Fatal(err)
		}

		switch args[0] {
		case "list":
			fmt.Println("Available versions:")
			displayAvailableVersions(r)
		case "latest":
			if err := getLatest(r); err != nil {
				logrus.Errorf("%s", err)
			} else {
				logrus.Infof("Success!")
			}
		default:
			if err := getSpecifiedVersion(r, args[0]); err != nil {
				logrus.Errorf("%s", err)
			} else {
				logrus.Infof("Success!")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// parseReleaseAPI populates the releases struct
func parseReleasesAPI() (releases, error) {
	r, err := http.Get("https://api.github.com/repos/eze-kiel/shaloc/releases")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var rel releases
	if err = json.Unmarshal(body, &rel); err != nil {
		return nil, err
	}

	return rel, nil
}

func displayAvailableVersions(r releases) {
	for i := 0; i < len(r); i++ {
		fmt.Printf("%s (%s)\n", aurora.BrightRed(r[i].TagName), aurora.Blue(r[i].PublishedAt))
	}
}

func getLatest(r releases) error {
	s := system{
		os:   runtime.GOOS,
		arch: runtime.GOARCH,
	}

	fullName := "shaloc_" + s.os + "_" + s.arch

	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	for archNum := 0; archNum < len(r[0].Assets); archNum++ {
		if r[0].Assets[archNum].Name == fullName {

			logrus.Info("Downloading shaloc:latest...")
			if err := download(binPath+"/shaloc-tmp", "https://github.com/eze-kiel/shaloc/releases/download/"+r[0].TagName+"/"+fullName); err != nil {
				return err
			}

			logrus.Infof("Installing shaloc:latest...")

			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Start()

			if err := os.Chmod(binPath+"/shaloc-tmp", 0775); err != nil {
				return err
			}

			if err := os.Remove(binPath + "/shaloc"); err != nil {
				return err
			}

			if err := os.Rename(binPath+"/shaloc-tmp", binPath+"/shaloc"); err != nil {
				return err
			}

			s.Stop()

			return nil
		}
	}
	return nil
}

func getVersionsList(r releases) []string {
	var versions []string
	for i := 0; i < len(r); i++ {
		versions = append(versions, r[i].TagName)
	}

	return versions
}

func getSpecifiedVersion(r releases, version string) error {
	s := system{
		os:   runtime.GOOS,
		arch: runtime.GOARCH,
	}

	fullName := "shaloc_" + s.os + "_" + s.arch

	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	versionsList := getVersionsList(r)
	if stringInSlice(version, versionsList) {
		logrus.Infof("Downloading shaloc:%s...", version)
		if err := download(binPath+"/shaloc-tmp", "https://github.com/eze-kiel/shaloc/releases/download/"+version+"/"+fullName); err != nil {
			return err
		}

		logrus.Infof("Installing shaloc:%s...", version)

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Start()

		if err := os.Chmod(binPath+"/shaloc-tmp", 0775); err != nil {
			return err
		}

		if err := os.Remove(binPath + "/shaloc"); err != nil {
			return err
		}

		if err := os.Rename(binPath+"/shaloc-tmp", binPath+"/shaloc"); err != nil {
			return err
		}

		s.Stop()
		return nil
	}

	return fmt.Errorf("Version %s not found", version)
}

// stringInSlice checks if a string appears in a slice.
func stringInSlice(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

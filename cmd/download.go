package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/dantdj/terraform-manager/config"
	"github.com/dantdj/terraform-manager/downloader"
	"github.com/dantdj/terraform-manager/parsing"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the specified versions of Terraform",
	Long:  "Downloads the specified versions of Terraform to the app directory for usage",
	PreRun: func(cmd *cobra.Command, args []string) {
		config.InitializeConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, version := range args {
			path, err := downloadTerraformVersion(version)
			if err != nil {
				return err
			}

			config.AddVersionConfig(version, path)
		}
		return nil
	},
}

func downloadTerraformVersion(version string) (string, error) {
	// TODO: Handle this error
	directory, _ := os.UserCacheDir()
	directory = directory + "/tfm"

	zipDest := fmt.Sprintf("%s/terraform_%s_%s_%s.zip", directory, version, runtime.GOOS, runtime.GOARCH)
	exeDest := fmt.Sprintf("%s/terraform/%s", directory, version)
	shaDest := directory + "/shasums.txt"

	url := fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip",
		version, version, runtime.GOOS, runtime.GOARCH,
	)
	fmt.Printf("Downloading Terraform v%s\n", version)
	if err := downloader.DownloadFile(url, zipDest); err != nil {
		return "", err
	}

	fmt.Printf("Downloading Terraform v%s hashsums\n", version)
	shaHashUrl := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_SHA256SUMS",
		version, version)
	if err := downloader.DownloadFile(shaHashUrl, shaDest); err != nil {
		return "", err
	}

	fmt.Println("")

	expectedHash, err := getExpectedHash(shaDest, zipDest)
	if err != nil {
		fmt.Println(err)
	}

	valid, err := downloader.ValidateFileHash(zipDest, expectedHash)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", fmt.Errorf("hash for downloaded file did not match the expected hash")
	}

	if err := downloader.UnzipFile(zipDest, exeDest); err != nil {
		return "", err
	}

	if err := os.Remove(zipDest); err != nil {
		return "", err
	}

	if err := os.Remove(shaDest); err != nil {
		return "", err
	}

	return exeDest, nil
}

func getExpectedHash(shaListFilepath, zipPath string) (string, error) {
	file, err := os.Open(shaListFilepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	parsedShaList := parsing.ParseShaList(file)

	// Doing this to remove directory detail - there's better ways of doing this but hacking for now
	// Should do something like having this be a proper path so I can just grab the filename
	directory, _ := os.UserCacheDir()
	directory = directory + "/tfm/"
	zipFileName := strings.Replace(zipPath, directory, "", -1)
	fmt.Printf("Final zip file name: %s", zipFileName)

	return parsedShaList[zipFileName], nil
}

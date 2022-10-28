package cmd

import (
	"fmt"
	"os"
	"runtime"

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
	zipDest := fmt.Sprintf("terraform_%s_%s_%s.zip", version, runtime.GOOS, runtime.GOARCH)
	exeDest := "terraform"
	finalPath := fmt.Sprintf("terraform/terraform%s", version)
	shaDest := "shasums.txt"

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

	expectedHash, _ := getExpectedHash(shaDest, zipDest)

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

	if err := os.Rename(exeDest+"/terraformtmp", exeDest+"/terraform"+version); err != nil {
		return "", err
	}

	if err := os.Remove(zipDest); err != nil {
		return "", err
	}

	if err := os.Remove(shaDest); err != nil {
		return "", err
	}

	return finalPath, nil
}

func getExpectedHash(shaListFilepath, zipPath string) (string, error) {
	file, err := os.Open(shaListFilepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	parsedShaList := parsing.ParseShaList(file)

	return parsedShaList[zipPath], nil
}

package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/dantdj/terraform-manager/downloader"
	"github.com/dantdj/terraform-manager/parsing"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the configured versions of Terraform",
	Long:  "Downloads configured versions of Terraform to the app directory for usage",
	RunE: func(cmd *cobra.Command, args []string) error {
		version := "1.0.9"
		zipDest := fmt.Sprintf("terraform_%s_%s_%s.zip", version, runtime.GOOS, runtime.GOARCH)
		exeDest := "terraform"
		shaDest := "shasums.txt"

		url := fmt.Sprintf(
			"https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip",
			version, version, runtime.GOOS, runtime.GOARCH,
		)
		if err := downloader.DownloadFile(url, zipDest); err != nil {
			return err
		}

		shaHashUrl := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_SHA256SUMS",
			version, version)
		if err := downloader.DownloadFile(shaHashUrl, shaDest); err != nil {
			return err
		}

		expectedHash, _ := getExpectedHash(shaDest, zipDest)

		valid, err := downloader.ValidateZipHash(zipDest, expectedHash)
		if err != nil {
			return err
		}
		if !valid {
			return fmt.Errorf("Hash for downloaded file did not match the expected hash")
		}

		if err := downloader.UnzipFile(zipDest, exeDest); err != nil {
			return err
		}

		if err := os.Rename(exeDest+"/terraform", exeDest+"/terraform"+version); err != nil {
			return err
		}

		if err := os.Remove(zipDest); err != nil {
			return err
		}

		if err := os.Remove(shaDest); err != nil {
			return err
		}

		return nil
	},
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

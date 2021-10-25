package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/dantdj/terraform-manager/downloader"
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
		zipDest := "terraform.zip"
		exeDest := "terraform"

		url := fmt.Sprintf(
			"https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip",
			version, version, runtime.GOOS, runtime.GOARCH,
		)
		if err := downloader.DownloadZip(url, zipDest); err != nil {
			return err
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

		return nil
	},
}

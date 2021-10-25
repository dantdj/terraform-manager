package downloader

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(progress []byte) (int, error) {
	totalLength := len(progress)
	wc.Total += uint64(totalLength)
	wc.PrintProgress()
	return totalLength, nil
}

func (wc *WriteCounter) PrintProgress() {
	// Clear line and print new download progress
	fmt.Print("\r")
	fmt.Printf("Downloading... %s downloaded", humanize.Bytes(wc.Total))
}

func DownloadZip(url string, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}

	response, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer response.Body.Close()

	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(response.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// Print newline once download finishes to move off the download line
	fmt.Print("\n")

	out.Close()

	return nil
}

func UnzipFile(source, dest string) ([]string, error) {
	var filenames []string

	reader, err := zip.OpenReader(source)
	if err != nil {
		return filenames, err
	}
	defer reader.Close()

	for _, file := range reader.File {
		destPath := filepath.Join(dest, file.Name)

		// Protects against ZipSlip - https://snyk.io/research/zip-slip-vulnerability
		if !strings.HasPrefix(destPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", destPath)
		}

		filenames = append(filenames, destPath)

		if file.FileInfo().IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := file.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}

	return filenames, nil
}

func ValidateHash(fileHash, expectedHash string) bool {
	return true
}

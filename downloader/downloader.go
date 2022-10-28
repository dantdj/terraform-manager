package downloader

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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
	// Space on the end required because the last character doesn't always clear properly
	fmt.Printf("Downloading... %s downloaded ", humanize.Bytes(wc.Total))
}

func DownloadFile(url string, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(response.Body, counter)); err != nil {
		return err
	}
	// Print newline once download finishes to move off the download line
	fmt.Print("\n")

	return nil
}

func UnzipFile(source, dest string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		destPath := dest

		if file.FileInfo().IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		contents, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, contents)

		// Close before next iteration of loop - deferring would only close at
		// end of function
		outFile.Close()
		contents.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateFileHash(filepath, expectedHash string) (bool, error) {
	hash, err := generateFileHash(filepath)
	if err != nil {
		return false, err
	}

	return hash == expectedHash, err
}

func generateFileHash(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

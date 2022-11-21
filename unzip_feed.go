package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip zip file "source" to destination.
// Stolen from https://gosamples.dev/unzip-file/
func unzipSource(source, destination string, desiredFiles []string) error {
	// Open the zip file
	fmt.Printf("Unzipping %s...\n", source)
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		if arrayContains(desiredFiles, f.Name) || len(desiredFiles) <= 0 {
			err := unzipFile(f, destination)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// stolen from https://gosamples.dev/unzip-file/
func unzipFile(f *zip.File, destination string) error {
	fmt.Printf("  %s\n", f.Name)
	// Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

func arrayContains(array []string, str string) bool {
	for _, i := range array {
		if i == str {
			return true
		}
	}
	return false
}

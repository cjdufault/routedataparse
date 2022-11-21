package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Downloads the GTFS feed file from url "feedUrl" and saves it to pwd.
// Stolen from https://golangdocs.com/golang-download-files
func downloadFeed(feedUrl string) string {

	// Build fileName from fullPath
	fileURL, err := url.Parse(feedUrl)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	fmt.Printf("Downloading %s...\n", fileName)
	resp, err := client.Get(feedUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("  Downloaded %s - Filesize: %d\n", fileName, size)
	return fileName
}

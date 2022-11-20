package main

import (
	"log"
)

func main() {

	//routeDataDestination := os.Args[1]

	// Download feed file from Metro Transit's GTFS service
	feedUrl := "https://svc.metrotransit.org/mtgtfs/gtfs.zip"
	feedFile := downloadFeed(feedUrl)

	// Unzip the file
	err := unzipSource(feedFile, ".")
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
)

func main() {

	//routeDataDestination := os.Args[1]
	neededFilenames := []string{"trips.txt", "stops.txt", "stop_times.txt"}

	// Download feed file from Metro Transit's GTFS service
	feedUrl := "https://svc.metrotransit.org/mtgtfs/gtfs.zip"
	feedFile := downloadFeed(feedUrl)

	// Unzip the file
	err := unzipSource(feedFile, ".", neededFilenames)
	if err != nil {
		log.Fatal(err)
	}
}

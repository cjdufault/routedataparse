package main

import (
	"fmt"
	"log"
)

func main() {

	//routeDataDestination := os.Args[1]
	routesFileName := "routes.txt"
	tripsFileName := "trips.txt"
	shapesFileName := "shapes.txt"
	neededFilenames := []string{routesFileName, tripsFileName, shapesFileName}

	// Download feed file from Metro Transit's GTFS service
	feedUrl := "https://svc.metrotransit.org/mtgtfs/gtfs.zip"
	feedFile := downloadFeed(feedUrl)

	// Unzip the file
	err := unzipSource(feedFile, ".", neededFilenames)
	if err != nil {
		log.Fatal(err)
	}

	shapes := getRouteShapes(routesFileName, tripsFileName, shapesFileName)
	fmt.Print(shapes[0])
}

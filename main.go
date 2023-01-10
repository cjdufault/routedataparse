package main

import (
	"os"
)

func main() {

	// handle args, and set defaults (for now)
	var feedUrl string
	var routeDataDestination string
	var stopDataDestination string

	if len(os.Args) >= 2 {
		feedUrl = os.Args[1]
	} else {
		feedUrl = "https://svc.metrotransit.org/mtgtfs/gtfs.zip"
	}
	if len(os.Args) >= 3 {
		routeDataDestination = os.Args[2]
	} else {
		routeDataDestination = "./shapes"
	}
	if len(os.Args) >= 4 {
		stopDataDestination = os.Args[3]
	} else {
		stopDataDestination = "./stops"
	}

	routesFileName := "routes.txt"
	tripsFileName := "trips.txt"
	shapesFileName := "shapes.txt"
	stopsFileName := "stops.txt"
	neededFilenames := []string{routesFileName, tripsFileName, shapesFileName, stopsFileName}

	// Download feed file from Metro Transit's GTFS service
	feedFile := downloadFeed(feedUrl)

	// Unzip the file
	if err := unzipSource(feedFile, ".", neededFilenames); err != nil {
		panic(err)
	}

	// Parse route data and output an array of Shape, one per route
	shapes := getRouteShapes(routesFileName, tripsFileName, shapesFileName)
	exportShapesToJson(shapes, routeDataDestination)

	// Parse stops data and output a json document of Stops
	stops := getStops(stopsFileName)
	exportStopsToJson(stops, stopDataDestination)
}

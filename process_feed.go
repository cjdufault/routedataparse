package main

import (
	"os"

	"github.com/gocarina/gocsv"

	. "github.com/ahmetb/go-linq/v3"
)

type Route struct {
	Id int `csv:"route_id"`
}

type Trip struct {
	Id        string `csv:"trip_id"`
	RouteId   int    `csv:"route_id"`
	Direction string `csv:"direction"`
	ShapeId   int    `csv:"shape_id"`
}

type ShapePoint struct {
	Id       int     `csv:"shape_id"`
	Sequence int     `csv:"shape_pt_sequence"`
	Lat      float64 `csv:"shape_pt_lat"`
	Lon      float64 `csv:"shape_pt_lon"`
}

type Shape struct {
	Id, RouteId int
	Direction   string
	Points      []ShapePoint
}

type Stop struct {
	Id   int     `csv:"stop_id"`
	Name string  `csv:"stop_name"`
	Lat  float64 `csv:"stop_lat"`
	Long float64 `csv:"stop_lon"`
}

// Reads data from routes.txt, trips.txt, and shapes.txt.
// Parses data into an array of Shape, with one shape per route.
// It's possible that this doesn't correctly handle variations in routes across trips.
func getRouteShapes(routesFileName string, tripsFileName string, shapesFileName string) []Shape {
	routes := getRoutes(routesFileName)
	trips := getTrips(tripsFileName)
	shapePoints := getShapePoints(shapesFileName)

	shapes := []Shape{}
	for _, route := range routes { // iterate through all routes

		directions := []string{}
		From(trips).Where(func(trip interface{}) bool {
			return trip.(*Trip).RouteId == route.Id
		}).Select(func(trip interface{}) interface{} {
			return trip.(*Trip).Direction
		}).Distinct().ToSlice(&directions)

		for _, direction := range directions {
			// from first trip with matching RouteId and matching direction, get ShapeId
			// this is where we may get bitten by variations in routes across trips
			shapeId := From(trips).Where(func(trip interface{}) bool {
				return trip.(*Trip).RouteId == route.Id
			}).Where(func(trip interface{}) bool {
				return trip.(*Trip).Direction == direction
			}).Select(func(trip interface{}) interface{} {
				return trip.(*Trip).ShapeId
			}).First().(int)

			// build array of ShapePoint where Id matches the ShapeId we got above
			shapePointArray := []ShapePoint{}
			From(shapePoints).Where(func(shapePoint interface{}) bool {
				return shapePoint.(*ShapePoint).Id == shapeId // match by ShapeId
			}).OrderBy(func(shapePoint interface{}) interface{} {
				return shapePoint.(*ShapePoint).Sequence // sort by Sequence
			}).Select(func(shapePoint interface{}) interface{} {
				var point ShapePoint
				point.Id = shapePoint.(*ShapePoint).Id
				point.Sequence = shapePoint.(*ShapePoint).Sequence
				point.Lat = shapePoint.(*ShapePoint).Lat
				point.Lon = shapePoint.(*ShapePoint).Lon
				return point
			}).ToSlice(&shapePointArray)

			// assign the values we've retrieved to a Shape, and append to the Shape array
			var shape Shape
			shape.Id = shapeId
			shape.RouteId = route.Id
			shape.Direction = direction
			shape.Points = shapePointArray

			shapes = append(shapes, shape)
		}
	}
	return shapes
}

// Reads route data in from routes.txt
func getRoutes(fileName string) []*Route {
	file := openFile(fileName)

	routes := []*Route{}
	if err := gocsv.UnmarshalFile(file, &routes); err != nil {
		panic(err)
	}

	defer file.Close()
	return routes
}

// Reads trip data in from trips.txt
func getTrips(fileName string) []*Trip {
	file := openFile(fileName)

	trips := []*Trip{}
	if err := gocsv.UnmarshalFile(file, &trips); err != nil {
		panic(err)
	}

	defer file.Close()
	return trips
}

// Reads shape data in from shapes.txt
func getShapePoints(fileName string) []*ShapePoint {
	file := openFile(fileName)

	shapePoints := []*ShapePoint{}
	if err := gocsv.UnmarshalFile(file, &shapePoints); err != nil {
		panic(err)
	}

	defer file.Close()
	return shapePoints
}

func getStops(stopsFileName string) []Stop {
	file := openFile(stopsFileName)

	stops := []Stop{}
	if err := gocsv.UnmarshalFile(file, &stops); err != nil {
		panic(err)
	}

	defer file.Close()
	return stops
}

func openFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return file
}

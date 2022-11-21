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
	Id      string `csv:"trip_id"`
	RouteId int    `csv:"route_id"`
	ShapeId int    `csv:"shape_id"`
}

type ShapePoint struct {
	Id       int     `csv:"shape_id"`
	Sequence int     `csv:"shape_pt_sequence"`
	Lat      float64 `csv:"shape_pt_lat"`
	Lon      float64 `csv:"shape_pt_lon"`
}

type Shape struct {
	Id, RouteId int
	Points      []ShapePoint
}

func getRouteShapes(routesFileName string, tripsFileName string, shapesFileName string) []Shape {
	routes := getRoutes(routesFileName)
	trips := getTrips(tripsFileName)
	shapePoints := getShapePoints(shapesFileName)

	shapes := []Shape{}
	for _, route := range routes {
		shapeId := From(trips).Where(func(trip interface{}) bool {
			return trip.(*Trip).RouteId == route.Id
		}).Select(func(trip interface{}) interface{} {
			return trip.(*Trip).ShapeId
		}).First().(int)

		shapePointArray := []ShapePoint{}
		From(shapePoints).Where(func(shapePoint interface{}) bool {
			return shapePoint.(*ShapePoint).Id == shapeId
		}).OrderBy(func(shapePoint interface{}) interface{} {
			return shapePoint.(*ShapePoint).Sequence
		}).Select(func(shapePoint interface{}) interface{} {
			var point ShapePoint
			point.Id = shapePoint.(*ShapePoint).Id
			point.Sequence = shapePoint.(*ShapePoint).Sequence
			point.Lat = shapePoint.(*ShapePoint).Lat
			point.Lon = shapePoint.(*ShapePoint).Lon
			return point
		}).ToSlice(&shapePointArray)

		var shape Shape
		shape.Id = shapeId
		shape.RouteId = route.Id
		shape.Points = shapePointArray

		shapes = append(shapes, shape)
	}

	return shapes
}

func getRoutes(fileName string) []*Route {
	file := openFile(fileName)

	routes := []*Route{}
	if err := gocsv.UnmarshalFile(file, &routes); err != nil {
		panic(err)
	}

	defer file.Close()
	return routes
}

func getTrips(fileName string) []*Trip {
	file := openFile(fileName)

	trips := []*Trip{}
	if err := gocsv.UnmarshalFile(file, &trips); err != nil {
		panic(err)
	}

	defer file.Close()
	return trips
}

func getShapePoints(fileName string) []*ShapePoint {
	file := openFile(fileName)

	shapePoints := []*ShapePoint{}
	if err := gocsv.UnmarshalFile(file, &shapePoints); err != nil {
		panic(err)
	}

	defer file.Close()
	return shapePoints
}

func openFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return file
}

func unmarshal(file *os.File, out interface{}) {

}

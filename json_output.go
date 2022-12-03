package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func exportShapesToJson(shapes []Shape, destination string) {

	fmt.Printf("Writing shape data files to %s...\n", destination)

	// Get the absolute destination path
	destination, err := filepath.Abs(destination)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		panic(err)
	}

	for _, shape := range shapes {
		shapeJson, err := json.Marshal(shape)

		if err != nil {
			panic(err)
		}

		filePath := filepath.Join(destination, fmt.Sprintf("%d_%s.json", shape.RouteId, shape.Direction))
		writeToFile(shapeJson, filePath)
	}
}

func writeToFile(contents []byte, path string) {
	fmt.Printf("  %s\n", path)
	if err := os.WriteFile(path, contents, 0666); err != nil {
		panic(err)
	}
}

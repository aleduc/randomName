package main

import (
	"flag"

	"test_quantcast/internal/app"
)

func main() {
	var (
		fileName, dateToProcess string
	)
	flag.StringVar(&fileName, "f", "", "file name")
	flag.StringVar(&dateToProcess, "d", "", "date to process")
	flag.Parse()

	app.Start(fileName, dateToProcess)
}

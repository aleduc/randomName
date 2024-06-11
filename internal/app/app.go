package app

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"test_quantcast/internal/logger"
	"test_quantcast/internal/parser"
	"test_quantcast/internal/printer"
	"test_quantcast/internal/processor"
)

const (
	columnsCount = 2
)

func Start(fileName, date string) {
	l := logger.NewSimple()
	osFile, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		l.Error(fmt.Errorf("open file: %w", err))
		return
	}
	defer func() {
		// for reading, it is not even an
		closeErr := osFile.Close()
		if closeErr != nil {
			l.Warn(fmt.Sprintf("close file: %s", closeErr.Error()))
		}
	}()

	rowProcessor, err := processor.NewRowProcessor(date, parser.NewParse(columnsCount))
	if err != nil {
		l.Error(fmt.Errorf("new row processor: %w", err))
	}
	csvProcessor := processor.NewCSV(l, rowProcessor, printer.NewResult(l))
	csvProcessor.ProcessFile(csv.NewReader(osFile))
}

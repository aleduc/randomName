package processor

import (
	"errors"
	"fmt"
	"io"

	"test_quantcast/internal/entity"
	"test_quantcast/internal/logger"
)

//go:generate go run github.com/golang/mock/mockgen --source=csv.go --destination=csv_mock.go --package=processor

type Reader interface {
	Read() (record []string, err error)
}

type RowProcessor interface {
	Process(input []string) error
	GetResult() entity.Result
}

type Printer interface {
	Print(result entity.Result)
}

type CSV struct {
	logger       logger.Logger
	rowProcessor RowProcessor
	printer      Printer
}

func NewCSV(logger logger.Logger, rowProcessor RowProcessor, printer Printer) *CSV {
	return &CSV{logger: logger, rowProcessor: rowProcessor, printer: printer}
}

func (c CSV) ProcessFile(csvReader Reader) {
	var n int
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			c.logger.Error(fmt.Errorf("csv read error: %w", err))
			return
		}
		if n == 0 {
			n++
			continue
		}

		err = c.rowProcessor.Process(row)
		if err != nil {
			if errors.Is(err, ErrNothingToProcess) {
				break
			}
			c.logger.Error(fmt.Errorf("process row: %w", err))
			return
		}

	}

	c.printer.Print(c.rowProcessor.GetResult())
}

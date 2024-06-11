package parser

import (
	"fmt"
	"time"

	"test_quantcast/internal/entity"
)

const (
	timestampLayout = "2006-01-02T15:04:05-07:00"
)

type Parse struct {
	columnsCount int
}

func NewParse(fieldsCount int) *Parse {
	return &Parse{columnsCount: fieldsCount}
}

func (p Parse) GetFileRow(input []string) (entity.FileRow, error) {
	if len(input) != p.columnsCount {
		return entity.FileRow{}, fmt.Errorf("to many columns. expected: %v, actual: %v", p.columnsCount, len(input))
	}

	timestampCol, err := time.Parse(timestampLayout, input[1])
	if err != nil {
		return entity.FileRow{}, fmt.Errorf("parse second column: %w", err)
	}
	return entity.FileRow{
		Cookie:    input[0],
		TimeStamp: timestampCol.UTC(),
	}, nil
}

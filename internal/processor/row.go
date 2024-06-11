package processor

import (
	"errors"
	"fmt"
	"time"

	"test_quantcast/internal/entity"
)

//go:generate go run github.com/golang/mock/mockgen --source=row.go --destination=row_mock.go --package=processor

type Parser interface {
	GetFileRow(input []string) (entity.FileRow, error)
}

const (
	calculationDateFormat = "2006-01-02"
)

var (
	ErrNothingToProcess = errors.New("nothing to process")
)

type Row struct {
	calculationDate time.Time
	parser          Parser
	max             int
	mapCnt          map[string]int
}

func NewRowProcessor(date string, parser Parser) (*Row, error) {
	t, err := time.Parse(calculationDateFormat, date)
	if err != nil {
		return nil, fmt.Errorf("parse calculation date: %w", err)
	}
	return &Row{calculationDate: t, parser: parser, mapCnt: make(map[string]int)}, nil
}

func (r *Row) Process(input []string) error {
	parsedRow, err := r.parser.GetFileRow(input)
	if err != nil {
		return fmt.Errorf("get file row: %w", err)
	}

	if parsedRow.TimeStamp.Before(r.calculationDate) {
		return nil
	}
	if parsedRow.TimeStamp.Truncate(time.Hour * 24).After(r.calculationDate) {
		return ErrNothingToProcess
	}

	r.mapCnt[parsedRow.Cookie]++
	if r.mapCnt[parsedRow.Cookie] > r.max {
		r.max = r.mapCnt[parsedRow.Cookie]
	}
	return nil
}

func (r *Row) GetResult() (res entity.Result) {
	res.TopCookies = make([]string, 0)
	for cookie, v := range r.mapCnt {
		if v == r.max {
			res.TopCookies = append(res.TopCookies, cookie)
		}
	}
	return
}

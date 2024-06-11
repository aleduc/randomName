package parser

import (
	"testing"
	"time"

	"test_quantcast/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestNewParse(t *testing.T) {
	tests := []struct {
		name         string
		fieldsCount  int
		expectedCols int
	}{
		{
			name:         "initialize with 2 columns",
			fieldsCount:  2,
			expectedCols: 2,
		},
		{
			name:         "initialize with 0 columns",
			fieldsCount:  0,
			expectedCols: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParse(tt.fieldsCount)
			assert.Equal(t, tt.expectedCols, parser.columnsCount)
		})
	}
}

func TestParse_GetFileRow(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		columnsCount  int
		expected      entity.FileRow
		expectedError string
	}{
		{
			name:         "valid input",
			input:        []string{"cookie1", "2023-01-01T12:00:00-05:00"},
			columnsCount: 2,
			expected: entity.FileRow{
				Cookie:    "cookie1",
				TimeStamp: time.Date(2023, 01, 01, 17, 00, 00, 0, time.UTC),
			},
			expectedError: "",
		},
		{
			name:          "too many columns",
			input:         []string{"cookie1", "2023-01-01T12:00:00-05:00", "extra"},
			columnsCount:  2,
			expected:      entity.FileRow{},
			expectedError: "to many columns. expected: 2, actual: 3",
		},
		{
			name:          "invalid timestamp format",
			input:         []string{"cookie1", "invalid-timestamp"},
			columnsCount:  2,
			expected:      entity.FileRow{},
			expectedError: "parse second column: parsing time \"invalid-timestamp\" as \"2006-01-02T15:04:05-07:00\": cannot parse \"invalid-timestamp\" as \"2006\"",
		},
		{
			name:          "not enough columns",
			input:         []string{"cookie1"},
			columnsCount:  2,
			expected:      entity.FileRow{},
			expectedError: "to many columns. expected: 2, actual: 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ao := assert.New(t)
			parser := NewParse(tt.columnsCount)
			result, err := parser.GetFileRow(tt.input)
			if tt.expectedError == "" {
				ao.NoError(err)
			} else {
				ao.EqualError(err, tt.expectedError)
			}
			ao.Equal(tt.expected, result)
		})
	}
}

package processor

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"test_quantcast/internal/entity"
)

func TestNewRowProcessor(t *testing.T) {
	tests := []struct {
		name         string
		date         string
		expectedDate time.Time
		expectErr    bool
	}{
		{
			name:         "valid date",
			date:         "2023-01-01",
			expectedDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expectErr:    false,
		},
		{
			name:      "invalid date",
			date:      "invalid-date",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockParser := NewMockParser(ctrl)

			rowProcessor, err := NewRowProcessor(tt.date, mockParser)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDate, rowProcessor.calculationDate)
				assert.Equal(t, mockParser, rowProcessor.parser)
				assert.NotNil(t, rowProcessor.mapCnt)
			}
		})
	}
}

func TestRow_Process(t *testing.T) {

	tests := []struct {
		name        string
		date        string
		input       []string
		expectErr   error
		expectMocks func(m *MockParser)
	}{
		{
			name:      "valid row within date range",
			date:      "2023-01-01",
			input:     []string{"cookie1", "2023-01-01T10:00:00Z"},
			expectErr: nil,
			expectMocks: func(m *MockParser) {
				m.EXPECT().GetFileRow([]string{"cookie1", "2023-01-01T10:00:00Z"}).Return(entity.FileRow{
					Cookie:    "cookie1",
					TimeStamp: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				}, nil).Times(1)
			},
		},
		{
			name:      "row before calculation date",
			date:      "2023-01-01",
			input:     []string{"cookie1", "2022-12-31T23:59:59Z"},
			expectErr: nil,
			expectMocks: func(m *MockParser) {
				m.EXPECT().GetFileRow([]string{"cookie1", "2022-12-31T23:59:59Z"}).Return(entity.FileRow{
					Cookie:    "cookie1",
					TimeStamp: time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC),
				}, nil).Times(1)
			},
		},
		{
			name:      "row after calculation date",
			date:      "2023-01-01",
			input:     []string{"cookie1", "2023-01-02T00:00:01Z"},
			expectErr: ErrNothingToProcess,
			expectMocks: func(m *MockParser) {
				m.EXPECT().GetFileRow([]string{"cookie1", "2023-01-02T00:00:01Z"}).Return(entity.FileRow{
					Cookie:    "cookie1",
					TimeStamp: time.Date(2023, 1, 2, 0, 0, 1, 0, time.UTC),
				}, nil).Times(1)
			},
		},
		{
			name:      "parser error",
			date:      "2023-01-01",
			input:     []string{"cookie1", "invalid-timestamp"},
			expectErr: fmt.Errorf("get file row: parsing time \"invalid-timestamp\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid-timestamp\" as \"2006\""),
			expectMocks: func(m *MockParser) {
				m.EXPECT().GetFileRow([]string{"cookie1", "invalid-timestamp"}).Return(entity.FileRow{}, fmt.Errorf("parsing time \"invalid-timestamp\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid-timestamp\" as \"2006\"")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ao := assert.New(t)
			ctrl := gomock.NewController(t)
			mockParser := NewMockParser(ctrl)

			rowProcessor, err := NewRowProcessor(tt.date, mockParser)
			ao.NoError(err)

			tt.expectMocks(mockParser)

			err = rowProcessor.Process(tt.input)
			if tt.expectErr != nil {
				ao.EqualError(err, tt.expectErr.Error())
			} else {
				ao.NoError(err)
			}
		})
	}
}

func TestRow_GetResult(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		processed map[string]int
		max       int
		expect    entity.Result
	}{
		{
			name: "single cookie",
			date: "2023-01-01",
			processed: map[string]int{
				"cookie1": 1,
			},
			max:    1,
			expect: entity.Result{TopCookies: []string{"cookie1"}},
		},
		{
			name: "multiple cookies with same count",
			date: "2023-01-01",
			processed: map[string]int{
				"cookie1": 2,
				"cookie2": 2,
			},
			max:    2,
			expect: entity.Result{TopCookies: []string{"cookie1", "cookie2"}},
		},
		{
			name: "one cookie with higher count",
			date: "2023-01-01",
			processed: map[string]int{
				"cookie1": 3,
				"cookie2": 2,
			},
			max:    3,
			expect: entity.Result{TopCookies: []string{"cookie1"}},
		},
		{
			name:      "no cookies",
			date:      "2023-01-01",
			processed: map[string]int{},
			max:       0,
			expect:    entity.Result{TopCookies: []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowProcessor := &Row{
				calculationDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				mapCnt:          tt.processed,
				max:             tt.max,
			}

			result := rowProcessor.GetResult()
			assert.Equal(t, tt.expect, result)
		})
	}
}

package processor

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"test_quantcast/internal/entity"
	"test_quantcast/internal/logger"
)

func TestNewCSV(t *testing.T) {
	ctrl := gomock.NewController(t)
	ao := assert.New(t)

	mockLogger := logger.NewMockLogger(ctrl)
	mockRowProcessor := NewMockRowProcessor(ctrl)
	mockPrinter := NewMockPrinter(ctrl)

	csvProcessor := NewCSV(mockLogger, mockRowProcessor, mockPrinter)

	ao.NotNil(csvProcessor)
	ao.Equal(mockLogger, csvProcessor.logger)
	ao.Equal(mockRowProcessor, csvProcessor.rowProcessor)
	ao.Equal(mockPrinter, csvProcessor.printer)

}
func TestCSV_ProcessFile(t *testing.T) {
	type mocks struct {
		reader       *MockReader
		rowProcessor *MockRowProcessor
		logger       *logger.MockLogger
		printer      *MockPrinter
	}
	tests := []struct {
		name          string
		records       [][]string
		processErrors []error
		expectPrint   entity.Result
		expectMocks   func(mocks)
	}{
		{
			name: "successful processing",
			records: [][]string{
				{"header1", "header2"},
				{"data1", "data2"},
			},
			processErrors: []error{nil, nil},
			expectPrint: entity.Result{
				TopCookies: []string{"cookie1"},
			},
			expectMocks: func(m mocks) {
				m.reader.EXPECT().Read().Return([]string{"header1", "header2"}, nil).Times(1)
				m.reader.EXPECT().Read().Return([]string{"data1", "data2"}, nil).Times(1)
				m.reader.EXPECT().Read().Return(nil, io.EOF).Times(1)
				m.rowProcessor.EXPECT().Process([]string{"data1", "data2"}).Return(nil).Times(1)
				m.rowProcessor.EXPECT().GetResult().Return(entity.Result{
					TopCookies: []string{"cookie1"},
				}).Times(1)
				m.printer.EXPECT().Print(entity.Result{
					TopCookies: []string{"cookie1"},
				}).Times(1)
			},
		},
		{
			name: "read error",
			records: [][]string{
				{"header1", "header2"},
			},
			processErrors: []error{io.EOF},
			expectPrint:   entity.Result{},
			expectMocks: func(m mocks) {
				m.reader.EXPECT().Read().Return([]string{"header1", "header2"}, nil).Times(1)
				m.reader.EXPECT().Read().Return(nil, errors.New("some err")).Times(1)
				m.logger.EXPECT().Error(fmt.Errorf("csv read error: %w", errors.New("some err"))).Times(1)
			},
		},
		{
			name: "process error",
			records: [][]string{
				{"header1", "header2"},
				{"data1", "data2"},
			},
			processErrors: []error{nil, errors.New("process error")},
			expectPrint:   entity.Result{},
			expectMocks: func(m mocks) {
				m.reader.EXPECT().Read().Return([]string{"header1", "header2"}, nil).Times(1)
				m.reader.EXPECT().Read().Return([]string{"data1", "data2"}, nil).Times(1)
				m.rowProcessor.EXPECT().Process([]string{"data1", "data2"}).Return(errors.New("process error")).Times(1)
				m.logger.EXPECT().Error(fmt.Errorf("process row: %w", errors.New("process error"))).Times(1)
			},
		},
		{
			name: "nothing to process error",
			records: [][]string{
				{"header1", "header2"},
				{"data1", "data2"},
			},
			processErrors: []error{nil, ErrNothingToProcess},
			expectPrint:   entity.Result{},
			expectMocks: func(m mocks) {
				m.reader.EXPECT().Read().Return([]string{"header1", "header2"}, nil).Times(1)
				m.reader.EXPECT().Read().Return([]string{"data1", "data2"}, nil).Times(1)
				m.rowProcessor.EXPECT().Process([]string{"data1", "data2"}).Return(ErrNothingToProcess).Times(1)
				m.rowProcessor.EXPECT().GetResult().Return(entity.Result{
					TopCookies: []string{},
				}).Times(1)
				m.printer.EXPECT().Print(entity.Result{
					TopCookies: []string{},
				}).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mocks{
				reader:       NewMockReader(ctrl),
				rowProcessor: NewMockRowProcessor(ctrl),
				logger:       logger.NewMockLogger(ctrl),
				printer:      NewMockPrinter(ctrl),
			}

			csvProcessor := NewCSV(m.logger, m.rowProcessor, m.printer)

			tt.expectMocks(m)

			csvProcessor.ProcessFile(m.reader)
		})
	}
}

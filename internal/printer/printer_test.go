package printer

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"test_quantcast/internal/entity"
	"test_quantcast/internal/logger"
)

func TestNewResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockLogger := logger.NewMockLogger(ctrl)
	result := NewResult(mockLogger)
	assert.NotNil(t, result)
	assert.Equal(t, mockLogger, result.l)

}

func TestResult_Print(t *testing.T) {
	tests := []struct {
		name   string
		result entity.Result
		prints []string
	}{
		{
			name: "print single cookie",
			result: entity.Result{
				TopCookies: []string{"cookie1"},
			},
			prints: []string{"cookie1"},
		},
		{
			name: "print multiple cookies",
			result: entity.Result{
				TopCookies: []string{"cookie1", "cookie2", "cookie3"},
			},
			prints: []string{"cookie1", "cookie2", "cookie3"},
		},
		{
			name: "print no cookies",
			result: entity.Result{
				TopCookies: []string{},
			},
			prints: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockLogger := logger.NewMockLogger(ctrl)
			result := NewResult(mockLogger)

			for _, p := range tt.prints {
				mockLogger.EXPECT().Println(p)
			}

			result.Print(tt.result)
		})
	}
}

package printer

import (
	"test_quantcast/internal/entity"
	"test_quantcast/internal/logger"
)

type Result struct {
	l logger.Logger
}

func NewResult(l logger.Logger) *Result {
	return &Result{l: l}
}

func (r Result) Print(result entity.Result) {
	for _, v := range result.TopCookies {
		r.l.Println(v)
	}
}

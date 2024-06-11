package logger

import (
	"fmt"
)

//go:generate go run github.com/golang/mock/mockgen --source=logger.go --destination=logger_mock.go --package=logger

type Logger interface {
	Println(input string)
	Error(args ...interface{})
	Warn(input string)
}

type Simple struct {
}

func NewSimple() *Simple {
	return &Simple{}
}

func (s *Simple) Println(input string) {
	fmt.Println(input)
}

func (s *Simple) Error(args ...interface{}) {
	fmt.Println(args...)
}

func (s *Simple) Warn(input string) {
	fmt.Println(input)
}

package ploto

import (
	"context"
	"fmt"
	"os"
)

type LoggerInterface interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	WithContext(ctx context.Context) LoggerInterface
}

type DefaultLogger struct {
}

func (l DefaultLogger) Debug(format string, v ...interface{}) {
	// fmt.FPrint(fmt.Sprintf(format, v...))
	fmt.Fprintln(os.Stdout, fmt.Sprintf(format, v...))
}

func (l DefaultLogger) Info(format string, v ...interface{}) {
	fmt.Fprintln(os.Stdout, fmt.Sprintf(format, v...))
}

func (l DefaultLogger) Warn(format string, v ...interface{}) {
	fmt.Fprintln(os.Stdout, fmt.Sprintf(format, v...))
}

func (l DefaultLogger) Error(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
}

func (l DefaultLogger) WithContext(ctx context.Context) LoggerInterface {
	return l
}

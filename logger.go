package ploto

import (
	"context"
	"fmt"
	"os"
)

type LoggerInterface interface {
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
}

type DefaultLogger struct {
}

func (l DefaultLogger) Debug(ctx context.Context, format string, v ...interface{}) {
	// fmt.FPrint(fmt.Sprintf(format, v...))
	fmt.Fprintln(os.Stdout, fmt.Sprintf(format, v...))
}

func (l DefaultLogger) Info(ctx context.Context, format string, v ...interface{}) {
	fmt.Fprintln(os.Stdout, fmt.Sprintf(format, v...))
}

func (l DefaultLogger) Warn(ctx context.Context, format string, v ...interface{}) {
	fmt.Fprintln(os.Stdout, fmt.Sprintf(format, v...))
}

func (l DefaultLogger) Error(ctx context.Context, format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
}

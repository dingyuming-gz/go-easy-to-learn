package logger

import (
	"fmt"
	"io"
	"log"
)

// LogLevel 定义了日志级别类型。
type LogLevel int

// 日志级别枚举值。
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger 封装了日志输出操作。
type Logger struct {
	logger *log.Logger
}

// NewLogger 创建一个新的日志实例。
func NewLogger(out io.Writer) *Logger {
	return &Logger{
		logger: log.New(out, "", log.Ldate|log.Ltime),
	}
}

// Logf 输出指定级别的日志。
func (l *Logger) Logf(format string, v ...interface{}) {
	l.logger.Printf(format+"\n", v...)
}

// SetLogLevel 设置日志输出的最小级别。
func (l *Logger) SetLogLevel(level LogLevel) {
	switch level {
	case DebugLevel:
		l.logger.SetPrefix("[DEBUG] ")
	case InfoLevel:
		l.logger.SetPrefix("[INFO] ")
	case WarnLevel:
		l.logger.SetPrefix("[WARN] ")
	case ErrorLevel:
		l.logger.SetPrefix("[ERROR] ")
	default:
		panic(fmt.Sprintf("unknown log level %d", level))
	}
}

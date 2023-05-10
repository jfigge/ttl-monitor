package common

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	LogLevelTrace = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type Logger struct {
	level LogLevel
}

var defaultLogger = NewLogger(LogLevelInfo)

func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level: level,
	}
}

func (l *Logger) write(level LogLevel, message string, args ...any) {
	if l.level <= level {
		fmt.Printf(message, args...)
	}
}

func (l *Logger) Trace(message string) {
	l.write(LogLevelTrace, "%s\n", message)
}
func (l *Logger) Debug(message string) {
	l.write(LogLevelDebug, "%s\n", message)
}
func (l *Logger) Info(message string) {
	l.write(LogLevelInfo, "%s\n", message)
}
func (l *Logger) Warn(message string) {
	l.write(LogLevelWarn, "%s\n", message)
}
func (l *Logger) Error(message string) error {
	l.write(LogLevelError, "%s\n", message)
	return fmt.Errorf(message)
}
func (l *Logger) Fatal(message string) {
	l.write(LogLevelFatal, "%s\n", message)
	os.Exit(1)
}

func (l *Logger) Tracef(message string, args ...any) {
	l.write(LogLevelTrace, message, args...)
}
func (l *Logger) Debugf(message string, args ...any) {
	l.write(LogLevelDebug, message, args...)
}
func (l *Logger) Infof(message string, args ...any) {
	l.write(LogLevelInfo, message, args...)
}
func (l *Logger) Warnf(message string, args ...any) {
	l.write(LogLevelWarn, message, args...)
}
func (l *Logger) Errorf(message string, args ...any) error {
	l.write(LogLevelError, message, args...)
	return fmt.Errorf(message, args...)
}
func (l *Logger) Fatalf(message string, args ...any) {
	l.write(LogLevelFatal, message, args...)
	os.Exit(1)
}

func Trace(message string) {
	defaultLogger.write(LogLevelTrace, "%s\n", message)
}
func Debug(message string) {
	defaultLogger.write(LogLevelDebug, "%s\n", message)
}
func Info(message string) {
	defaultLogger.write(LogLevelInfo, "%s\n", message)
}
func Warn(message string) {
	defaultLogger.write(LogLevelWarn, "%s\n", message)
}
func Error(message string) error {
	defaultLogger.write(LogLevelError, "%s\n", message)
	return fmt.Errorf(message)
}
func Fatal(message string) {
	defaultLogger.write(LogLevelFatal, "%s\n", message)
	os.Exit(1)
}

func Tracef(message string, args ...any) {
	defaultLogger.write(LogLevelTrace, message, args...)
}
func Debugf(message string, args ...any) {
	defaultLogger.write(LogLevelDebug, message, args...)
}
func Infof(message string, args ...any) {
	defaultLogger.write(LogLevelInfo, message, args...)
}
func Warnf(message string, args ...any) {
	defaultLogger.write(LogLevelWarn, message, args...)
}
func Errorf(message string, args ...any) error {
	defaultLogger.write(LogLevelError, message, args...)
	return fmt.Errorf(message, args...)
}
func Fatalf(message string, args ...any) {
	defaultLogger.write(LogLevelFatal, message, args...)
	os.Exit(1)
}

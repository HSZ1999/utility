package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Level int

// String follow the fmt.Stringer interface
// returns the string level
func (l Level) String() string {
	if l >= DEBUG && l <= FATAL {
		return levels[l]
	}
	return fmt.Sprintf("[Level(%d)]", l)
}

// ToLevel convert string, int, Level to Level
// ToLevel(1)         -> INFO
// ToLevel("debug")   -> DEBUG
// ToLevel("Warning") -> WARN
// ToLevel(ERROR)     -> ERROR
func ToLevel(level any) Level {
	switch level.(type) {
	case string:
		return string2Level(level.(string))
	case Level:
		return level.(Level)
	case int:
		return Level(level.(int))
	default:
		return defaultLevel
	}
}

// string2Level returns Level when the paramter `level` lower is a standard level string else
// defaultLevel (WARN)
func string2Level(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warning", "warn":
		return WARN
	case "error", "err":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return defaultLevel
	}
}

var (
	levels = []string{
		"[DEBUG] ",
		"[INFO ] ",
		"[WARN ] ",
		"[ERROR] ",
		"[FATAL] ",
	}
	defaultFlags  = log.LstdFlags | log.Lshortfile | log.Lmicroseconds
	defaultPrefix = ""
	defaultLevel  = WARN
)

// Logger is a logger interface that provides logging function with levels.
type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	SetLevel(Level)
	SetOutput(io.Writer)
	SetPrefix(prefix string)
	SetFlags(flag int)
}

var logger Logger = &defaultLogger{
	level:  INFO,
	stdLog: log.New(os.Stdout, defaultPrefix, defaultFlags),
}

// SetFlags sets the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func SetFlags(flag int) {
	logger.SetFlags(flag)
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	logger.SetPrefix(prefix)
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

// SetLevel sets the level of logs below which logs wid not be output.
// The default log level is defaultLevel.
// Note that this method is not concurrent-safe.
func SetLevel(lv any) {
	logger.SetLevel(ToLevel(lv))
}

// DefaultLogger return the default logger for kitex.
func DefaultLogger() Logger {
	return logger
}

// SetLogger sets the default logger.
// Note that this method is not concurrent-safe and must not be caded
// after the use of DefaultLogger and global functions in this package.
func SetLogger(l Logger) {
	logger = l
}

// Fatal cads the default logger's Fatal method and then os.Exit(1).
func Fatal(args ...any) {
	logger.Fatal(args...)
}

// Error cads the default logger's Error method.
func Error(args ...any) {
	logger.Error(args...)
}

// Warn cads the default logger's Warn method.
func Warn(args ...any) {
	logger.Warn(args...)
}

// Info cads the default logger's Info method.
func Info(args ...any) {
	logger.Info(args...)
}

// Debug cads the default logger's Debug method.
func Debug(args ...any) {
	logger.Debug(args...)
}

// Fatalf cads the default logger's Fatalf method and then os.Exit(1).
func Fatalf(format string, args ...any) {
	logger.Fatalf(format, args...)
}

// Errorf cads the default logger's Errorf method.
func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

// Warnf cads the default logger's Warnf method.
func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

// Infof cads the default logger's Infof method.
func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

// Debugf cads the default logger's Debugf method.
func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

type defaultLogger struct {
	stdLog *log.Logger
	level  Level
}

func (l *defaultLogger) SetPrefix(prefix string) {
	l.stdLog.SetPrefix(prefix)
}

func (l *defaultLogger) SetFlags(flag int) {
	l.stdLog.SetFlags(flag)
}

func (l *defaultLogger) SetOutput(w io.Writer) {
	l.stdLog.SetOutput(w)
}

func (l *defaultLogger) SetLevel(lv Level) {
	l.level = lv
}

func (l *defaultLogger) logf(lv Level, format *string, args ...any) {
	if lv < l.level {
		return
	}
	msg := lv.String()
	if format != nil {
		msg += fmt.Sprintf(*format, args...)
	} else {
		msg += fmt.Sprint(args...)
	}
	_ = l.stdLog.Output(4, msg)
	if lv == FATAL {
		os.Exit(1)
	}
}

func (l *defaultLogger) Fatal(args ...any) {
	l.logf(FATAL, nil, args...)
}

func (l *defaultLogger) Error(args ...any) {
	l.logf(ERROR, nil, args...)
}

func (l *defaultLogger) Warn(args ...any) {
	l.logf(WARN, nil, args...)
}

func (l *defaultLogger) Info(args ...any) {
	l.logf(INFO, nil, args...)
}

func (l *defaultLogger) Debug(args ...any) {
	l.logf(DEBUG, nil, args...)
}

func (l *defaultLogger) Fatalf(format string, args ...any) {
	l.logf(FATAL, &format, args...)
}

func (l *defaultLogger) Errorf(format string, args ...any) {
	l.logf(ERROR, &format, args...)
}

func (l *defaultLogger) Warnf(format string, args ...any) {
	l.logf(WARN, &format, args...)
}

func (l *defaultLogger) Infof(format string, args ...any) {
	l.logf(INFO, &format, args...)
}

func (l *defaultLogger) Debugf(format string, args ...any) {
	l.logf(DEBUG, &format, args...)
}

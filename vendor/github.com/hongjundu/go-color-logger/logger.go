package logger

import (
	"fmt"
	"github.com/mgutz/logxi/v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

var (
	stdLogger log.Logger
	errLogger log.Logger
)

func Init(logFileName string, logFileDir string, maxLogFileSize, maxLogFileBackups, maxLogFileAge int) {
	stdLoggerName := "std"
	errLoggerName := "err"

	// Init default logger
	stdLogger = log.NewLogger(os.Stdout, stdLoggerName)
	errLogger = log.NewLogger(os.Stderr, errLoggerName)

	if len(logFileDir) >= 0 {
		dirExists, _ := dirExists(logFileDir)
		if dirExists {
			maxSize := maxLogFileSize
			maxBackups := maxLogFileBackups
			maxAge := maxLogFileAge

			if maxSize <= 0 {
				maxSize = 100
			}

			if maxBackups <= 0 {
				maxBackups = 3
			}

			if maxAge <= 0 {
				maxAge = 30
			}

			stdOutputFile := &lumberjack.Logger{
				Filename:   filepath.Join(logFileDir, fmt.Sprintf("%s.log", logFileName)),
				MaxSize:    maxSize, // megabytes
				MaxBackups: maxBackups,
				MaxAge:     maxAge, //days
			}

			errOutputFile := &lumberjack.Logger{
				Filename:   filepath.Join(logFileDir, fmt.Sprintf("%s-error.log", logFileName)),
				MaxSize:    maxSize, // megabytes
				MaxBackups: maxBackups,
				MaxAge:     maxAge, //days
			}

			stdOutput := io.MultiWriter(stdOutputFile, os.Stdout)
			errOutput := io.MultiWriter(errOutputFile, os.Stderr)

			stdLogger = log.NewLogger(stdOutput, stdLoggerName)
			errLogger = log.NewLogger(errOutput, errLoggerName)
		} else {
			stdLogger.Warn("logFileDir does not exit", "dir", logFileDir)
		}
	} else {
		stdLogger.Info("logFileDir is empty")
	}

}

func dirExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err == nil {
		return true, nil
	} else {
		if os.IsNotExist(err) {
			return false, err
		} else {
			return true, err
		}
	}
}

func Trace(msg string, args ...interface{}) {
	stdLogger.Trace(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	stdLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	stdLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) error {
	stdLogger.Warn(msg, args...)
	return errLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) error {
	stdLogger.Error(msg, args...)
	return errLogger.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	stdLogger.Fatal(msg, args...)
}

func Log(level int, msg string, args []interface{}) {
	stdLogger.Log(level, msg, args)
}

func SetLevel(level int) {
	stdLogger.SetLevel(level)
}

func IsTrace() bool {
	return stdLogger.IsTrace()
}

func IsDebug() bool {
	return stdLogger.IsDebug()
}

func IsInfo() bool {
	return stdLogger.IsInfo()
}

func IsWarn() bool {
	return stdLogger.IsWarn()
}

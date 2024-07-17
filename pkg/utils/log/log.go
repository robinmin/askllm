package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"log/slog"

	"github.com/dusted-go/logging/prettylog"
)

func InitLogger(logpath string, app_id string, level string, verbose bool) (logFile *os.File) {
	// Set the custom logger as the default
	// var logFile *os.File
	lvl, err := ParseLevel(strings.ToUpper(level))
	if err == nil || verbose {
		lvl = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	if logpath != "" {
		// Create a custom JSON logger
		filename := fmt.Sprintf("%s/%s_%s.log", logpath, app_id, time.Now().Format("20060102"))

		if logFile == nil {
			err := ensureDirExists(filename)
			if err != nil {
				slog.Error((err.Error()))
				return nil
			}

			var err1 error
			logFile, err1 = os.Create(filename)
			if err1 != nil {
				slog.Error((err1.Error())) // Handle errors appropriately
				return nil
			}
			logger := slog.New(slog.NewJSONHandler(logFile, opts))
			slog.SetDefault(logger)
		}
	} else {
		// logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
		prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
			Level:       slog.LevelInfo,
			AddSource:   false,
			ReplaceAttr: nil,
		})
		logger := slog.New(prettyHandler)
		slog.SetDefault(logger)
	}
	return logFile
}

func CloseLogger(logFile *os.File) {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			fmt.Printf("Error closing log file: %v", err)
		}
		logFile = nil
	}
}

// ensureDirExists checks if the directory for the given file path exists.
// If it doesn't exist, it creates the directory.
// It returns an error if the directory couldn't be created.
func ensureDirExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		// fmt.Printf("Created directory: %s\n", dir)
	}
	return nil
}

func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}

func Debug(msg string, args ...interface{}) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	slog.Error(msg, args...)
}

func Debugf(msg string, args ...interface{}) {
	slog.Debug(fmt.Sprintf(msg, args...))
}

func Infof(msg string, args ...interface{}) {
	slog.Info(fmt.Sprintf(msg, args...))
}

func Warnf(msg string, args ...interface{}) {
	slog.Warn(fmt.Sprintf(msg, args...))
}

func Errorf(msg string, args ...interface{}) {
	slog.Error(fmt.Sprintf(msg, args...))
}

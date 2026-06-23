/*
Copyright 2026 CodeFuture Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package log provides slog-based structured logging for kube-agents.
package log

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Options configures the logger.
type Options struct {
	Level      string // debug, info, warn, error
	Format     string // text, json
	AddSource  bool   // include source file and line
	FileOutput bool   // write to file instead of stderr
	FilePath   string // log file path
	MaxSize    int    // max megabytes before rotation (default: 100)
	MaxBackups int    // max old files to retain (default: 10)
	MaxAge     int    // max days to retain (default: 30)
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Level:      "info",
		Format:     "text",
		AddSource:  false,
		FileOutput: false,
		FilePath:   "",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
	}
}

// Init initializes the global slog logger with the given options.
func Init(opts Options) {
	level := parseLevel(opts.Level)
	handlerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: opts.AddSource,
	}

	var w io.Writer = os.Stderr

	if opts.FileOutput && opts.FilePath != "" {
		w = &lumberjack.Logger{
			Filename:   opts.FilePath,
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   true,
		}
	}

	var handler slog.Handler
	switch opts.Format {
	case "json":
		handler = slog.NewJSONHandler(w, handlerOpts)
	default:
		handler = slog.NewTextHandler(w, handlerOpts)
	}

	slog.SetDefault(slog.New(handler))
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

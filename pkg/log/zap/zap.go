/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package zap

import (
	"errors"
	"strings"

	"adeia/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// levels is a map of supported log levels.
var levels = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"panic": zap.PanicLevel,
	"fatal": zap.FatalLevel,
}

// Logger represents a logger that can log messages.
type Logger struct {
	*zap.SugaredLogger
}

// New creates a new Logger with the specified conf.
func New(conf *config.LoggerConfig) (*Logger, error) {
	level, err := parseLevel(conf.Level)
	if err != nil {
		return nil, err
	}

	// TODO: switch to custom config
	// TODO: setup log rotation
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.OutputPaths = conf.Paths

	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{l.Sugar()}, nil
}

// parseLevel returns the appropriate zapcore.Level for the passed-in string.
func parseLevel(s string) (zapcore.Level, error) {
	if l, ok := levels[strings.ToLower(s)]; ok {
		return l, nil
	}

	return 0, errors.New("specified log level is not one of ['debug', 'info', 'warn', 'error', 'panic', 'fatal']")
}

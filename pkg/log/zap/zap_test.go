/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package zap

import (
	"testing"

	"adeia/internal/config"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestParseLevel(t *testing.T) {
	t.Run("parse when level is a valid log level", func(t *testing.T) {
		t.Parallel()
		want := zap.InfoLevel
		got, err := parseLevel("info")
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("return error when level is invalid", func(t *testing.T) {
		t.Parallel()
		_, err := parseLevel("info123")
		assert.Error(t, err, "should return error when string is invalid")
	})
}

func TestNew(t *testing.T) {
	t.Run("return error on invalid level", func(t *testing.T) {
		t.Parallel()
		c := &config.LoggerConfig{Level: "foobar123"}
		_, err := New(c)
		assert.Error(t, err)
	})

	t.Run("return error on invalid paths", func(t *testing.T) {
		t.Parallel()
		c := &config.LoggerConfig{
			Level: "debug",
			Paths: []string{`/`},
		}
		_, err := New(c)
		assert.Error(t, err)
	})

	t.Run("return logger on valid config", func(t *testing.T) {
		c := &config.LoggerConfig{
			Level: "debug",
			Paths: []string{"stderr"},
		}
		got, err := New(c)

		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.OutputPaths = []string{"stderr"}
		l, _ := cfg.Build(zap.AddCallerSkip(1))
		want := &Logger{l.Sugar()}

		assert.Nil(t, err)
		if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(zap.SugaredLogger{})) {
			t.Errorf("want: %v, got :%v", want, got)
		}
	})
}

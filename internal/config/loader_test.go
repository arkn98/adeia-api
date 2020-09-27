/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("load valid config", func(t *testing.T) {
		t.Parallel()
		b := bytes.NewBufferString(`
server:
  host: test
  port: 1234
logger:
  level: info
  paths:
    - foo
    - bar
`)
		want := Config{
			ServerConfig: ServerConfig{
				Port: 1234,
				Host: "test",
			},
			LoggerConfig: LoggerConfig{
				Level: "info",
				Paths: []string{"foo", "bar"},
			},
		}

		got, err := Load(b)
		assert.Nil(t, err)
		assert.Equal(t, want.ServerConfig.Port, got.ServerConfig.Port)
		assert.Equal(t, want.ServerConfig.Host, got.ServerConfig.Host)
		assert.Equal(t, want.LoggerConfig.Level, got.LoggerConfig.Level)
	})

	t.Run("return error on invalid yaml", func(t *testing.T) {
		t.Parallel()
		b := bytes.NewBufferString(`
@
`)
		_, err := Load(b)
		assert.Error(t, err)
	})

	t.Run("return error when unmarshalling fails", func(t *testing.T) {
		t.Parallel()
		b := bytes.NewBufferString(`
server:
  port: hello
`)
		_, err := Load(b)
		assert.Error(t, err)
	})

	t.Run("env should override config", func(t *testing.T) {
		want := "bar"
		key := "foobar"
		b := bytes.NewBufferString(`
server:
  jwt_secret: foo
`)
		overrides := map[string]string{
			"server.jwt_secret": key,
		}
		_ = os.Setenv(key, want)
		applyEnvOverrides = envOverride(overrides)
		got, err := Load(b)
		assert.Nil(t, err)
		assert.Equal(t, want, got.ServerConfig.JWTSecret)
	})

	t.Run("wrong env should not override config", func(t *testing.T) {
		want := "foo"
		key := "foobar"
		b := bytes.NewBufferString(`
server:
  jwt_secret: foo
`)
		overrides := map[string]string{
			"server.jwt_secret1": key,
		}

		_ = os.Setenv(key, want)
		defer os.Unsetenv(key)

		applyEnvOverrides = envOverride(overrides)
		got, err := Load(b)
		assert.Nil(t, err)
		assert.Equal(t, want, got.ServerConfig.JWTSecret)
	})
}

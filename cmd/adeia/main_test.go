/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("return fallback if key is not set", func(t *testing.T) {
		t.Parallel()

		want := "bar"
		key := "DUMMY_KEY"
		_ = os.Unsetenv(key)
		got := getEnv(key, "bar")
		assert.Equal(t, want, got)
	})

	t.Run("return value from env if key is set", func(t *testing.T) {
		t.Parallel()

		key := "DUMMY_KEY"
		want := "foo"
		_ = os.Setenv(key, want)
		defer func() {
			_ = os.Unsetenv(key)
		}()
		got := getEnv("DUMMY_KEY", "bar")
		assert.Equal(t, want, got)
	})
}

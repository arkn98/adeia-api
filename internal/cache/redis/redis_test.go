/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package redis

import (
	"strconv"
	"testing"

	"adeia/internal/config"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (r *Redis, mock *miniredis.Miniredis, c func()) {
	t.Parallel()
	mock, _ = miniredis.Run()
	c = func() {
		mock.Close()
	}

	port, _ := strconv.Atoi(mock.Port())
	conf := &config.CacheConfig{
		Network:  "tcp",
		Host:     mock.Host(),
		Port:     port,
		ConnSize: 10,
	}
	r, _ = New(conf)

	return
}

func TestNew(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		t.Parallel()
		c := &config.CacheConfig{
			Network:  "tcp",
			Host:     "256.256.256.256",
			Port:     5432,
			ConnSize: 10,
		}
		_, err := New(c)
		assert.Error(t, err)
	})
}

func TestRedis_Get(t *testing.T) {
	t.Run("return value when key exists", func(t *testing.T) {
		r, mock, c := setup(t)
		defer c()

		key := "foo"
		want := "bar"
		_ = mock.Set(key, want)
		var got string
		err := r.Get(&got, "foo")

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("error when value is not string", func(t *testing.T) {
		r, mock, c := setup(t)
		defer c()

		key := "foo"
		_, _ = mock.RPush(key, "bar", "baz")
		var got *struct {
			Test string
		}
		err := r.Get(got, key)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRedis_Set(t *testing.T) {
	t.Run("set value", func(t *testing.T) {
		r, mock, c := setup(t)
		defer c()

		key := "foo"
		want := "bar"
		err := r.Set(key, want)
		got, _ := mock.Get(key)

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})
}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package redis

import (
	"strconv"

	"adeia/internal/config"

	"github.com/mediocregopher/radix/v3"
)

// Redis represents the cache connection instance.
type Redis struct {
	*radix.Pool
}

// New creates a new cache connection instance.
func New(conf *config.CacheConfig) (*Redis, error) {
	// TODO: add cache auth
	p, err := radix.NewPool(
		conf.Network,
		conf.Host+":"+strconv.Itoa(conf.Port),
		conf.ConnSize,
	)
	if err != nil {
		return nil, err
	}

	return &Redis{p}, nil
}

// Get gets the value of the specified key.
func (r *Redis) Get(dest interface{}, key string) error {
	return r.Do(radix.Cmd(dest, "GET", key))
}

// Set sets the provided key:value pair.
func (r *Redis) Set(key string, value string) error {
	return r.Do(radix.Cmd(nil, "SET", key, value))
}

/*
// Delete deletes the list of keys.
func (r *Redis) Delete(keys ...string) error {
	return r.Do(radix.Cmd(nil, "DEL", keys...))
}

// SetWithExpiry sets the provided key:value pair with specified seconds of TTL.
func (r *Redis) SetWithExpiry(key, value string, seconds int) error {
	return r.Do(radix.Cmd(nil, "SET", key, value, "EX", strconv.Itoa(seconds)))
}

// Expire sets the expiry for a given key.
func (r *Redis) Expire(key string, seconds int) error {
	return r.Do(radix.Cmd(nil, "EXPIRE", key, strconv.Itoa(seconds)))
}

func buildKey(resource, id string, fields ...string) string {
	key := resource + ":" + id
	for _, field := range fields {
		key += ":" + field
	}
	return key
}
*/

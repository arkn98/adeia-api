/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cache

// Cache is an interface for all cache-related functions, that implementations
// must implement.
type Cache interface {
	Close() error
	Get(dest interface{}, key string) error
	Set(key string, value string) error
}

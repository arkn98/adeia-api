/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package log

// Logger is the interface for all the functions of a logger.
type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Sync() error
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
}

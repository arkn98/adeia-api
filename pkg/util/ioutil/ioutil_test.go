/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package ioutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errFailedToClose = errors.New("failed to close")

type MockCloser struct {
	shouldFail bool
	isClosed   bool
}

func (mc *MockCloser) Close() error {
	if mc.shouldFail {
		return errFailedToClose
	}
	mc.isClosed = true
	return nil
}

func TestCheckClose(t *testing.T) {
	t.Run("successful close + no previous error", func(t *testing.T) {
		t.Parallel()
		var err error = nil
		c := &MockCloser{}

		CheckCloseErr(c, &err)
		assert.Nil(t, err)
		assert.True(t, c.isClosed)
	})

	t.Run("successful close + previous error", func(t *testing.T) {
		t.Parallel()
		err := errors.New("error")
		c := &MockCloser{}

		CheckCloseErr(c, &err)
		assert.Error(t, err)
		assert.EqualError(t, err, "error")
		assert.True(t, c.isClosed)
	})

	t.Run("error while closing + no previous error", func(t *testing.T) {
		t.Parallel()
		var err error = nil
		c := &MockCloser{shouldFail: true}

		CheckCloseErr(c, &err)
		assert.Error(t, err)
		assert.EqualError(t, err, errFailedToClose.Error())
		assert.False(t, c.isClosed)
	})

	t.Run("error while closing + previous error", func(t *testing.T) {
		t.Parallel()
		err := errors.New("error")
		c := &MockCloser{shouldFail: true}

		CheckCloseErr(c, &err)
		assert.Error(t, err)
		assert.EqualError(t, err, "error")
		assert.False(t, c.isClosed)
	})
}

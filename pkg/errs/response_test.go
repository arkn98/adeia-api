/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseError_Msg(t *testing.T) {
	t.Parallel()
	r := ResponseError{}
	want := "foobar"
	r = r.Msg(want)
	assert.Equal(t, want, r.Message)
}

func TestResponseError_Msgf(t *testing.T) {
	t.Parallel()
	r := ResponseError{}
	want := fmt.Sprintf("foo %v %v %v", "bar", "baz", 10)
	r = r.Msgf("foo %v %v %v", "bar", "baz", 10)
	assert.Equal(t, want, r.Message)
}

func TestResponseError_AddValidationErr(t *testing.T) {
	t.Run("empty validation error map", func(t *testing.T) {
		t.Parallel()
		r := ResponseError{}
		assert.Empty(t, r.ValidationErrors)
	})

	t.Run("single validation error", func(t *testing.T) {
		t.Parallel()
		r := ResponseError{}
		r = r.AddValidationErr("foo", "bar")
		assert.Equal(t, map[string]string{"foo": "bar"}, r.ValidationErrors)
		assert.Len(t, r.ValidationErrors, 1)
	})

	t.Run("validation error with same keys override", func(t *testing.T) {
		t.Parallel()
		r := ResponseError{}
		r = r.AddValidationErr("foo", "val1")
		r = r.AddValidationErr("bar", "val2")
		r = r.AddValidationErr("foo", "val3")

		assert.Equal(t, map[string]string{"foo": "val3", "bar": "val2"}, r.ValidationErrors)
		assert.Len(t, r.ValidationErrors, 2)
	})
}

func TestResponseError_ValidationErr(t *testing.T) {
	t.Run("set empty map", func(t *testing.T) {
		t.Parallel()
		r := ResponseError{}
		r = r.ValidationErr(map[string]string{})
		assert.Empty(t, r.ValidationErrors)
	})

	t.Run("set validation error map", func(t *testing.T) {
		t.Parallel()
		r := ResponseError{}
		v := map[string]string{"foo": "val3", "bar": "val2"}
		r = r.ValidationErr(v)
		assert.Equal(t, v, r.ValidationErrors)
		assert.Len(t, r.ValidationErrors, 2)
	})
}

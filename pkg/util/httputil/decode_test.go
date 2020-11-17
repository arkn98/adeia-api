/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package httputil

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"adeia"

	"github.com/stretchr/testify/assert"
)

type MockWarner struct {
	bytes.Buffer
}

func (mw *MockWarner) Warnf(template string, args ...interface{}) {
	mw.Buffer.WriteString(fmt.Sprintf(template, args...))
}

func TestLogWriteErr(t *testing.T) {
	t.Run("write error", func(t *testing.T) {
		t.Parallel()
		m := &MockWarner{}
		e := errors.New("test error")
		LogWriteErr(m, e)
		want := fmt.Errorf("failed to write response: %v", e)
		assert.Equal(t, want.Error(), m.Buffer.String())
	})

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()
		m := &MockWarner{}
		LogWriteErr(m, nil)
		assert.Empty(t, m.Buffer)
	})
}

func TestIsReqMalformedErr(t *testing.T) {
	testcases := []struct {
		in   error
		want bool
		msg  string
	}{
		{adeia.ErrInvalidRequest, true, "req malformed error"},
		{adeia.ErrUnsupportedMediaType, true, "req malformed error"},
		{adeia.ErrValidationFailed, true, "req malformed error"},
		{errors.New("test"), false, "some other error"},
	}
	for _, tc := range testcases {
		t.Run(tc.msg, func(t *testing.T) {
			t.Parallel()
			got := isReqMalformedErr(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

type errResponseWriter struct {
	*httptest.ResponseRecorder
}

func (er *errResponseWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("test error")
}

func TestDecode(t *testing.T) {
	t.Run("request malformed error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		body := `
foo:
  bar
`
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		err := Decode(w, r, nil)

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "malformed request body")
		}
	})

	t.Run("successful decode", func(t *testing.T) {
		t.Parallel()

		type body struct {
			Foo string `json:"foo"`
			Baz int    `json:"baz"`
		}

		var got body
		want := body{"bar", 10}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`{"foo":"bar","baz":10}`))
		r.Header.Set("Content-Type", "application/json")
		err := Decode(w, r, &got)

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})
}

func TestDecodeJSONBody(t *testing.T) {
	setup := func(t *testing.T) *httptest.ResponseRecorder {
		t.Parallel()
		return httptest.NewRecorder()
	}

	t.Run("invalid Content-Type", func(t *testing.T) {
		w := setup(t)
		r := httptest.NewRequest(http.MethodGet, "/1", nil)
		r.Header.Set("Content-Type", "application/xml")
		err := decodeJSONBody(w, r, nil, 131072)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrUnsupportedMediaType.Error())
	})

	t.Run("syntax error in json", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo int `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`foo`))
		r.Header.Set("Content-Type", "application/json")
		err := decodeJSONBody(w, r, &got, 131072)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidRequest.Error())
	})

	t.Run("json validation failed", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo int `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`{"foo":["bar","baz"]}`))
		r.Header.Set("Content-Type", "application/json")
		err := decodeJSONBody(w, r, &got, 131072)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrValidationFailed.Error())
	})

	t.Run("unknown json fields", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo int `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`{"foo":1,"bar":"baz"}`))
		r.Header.Set("Content-Type", "application/json")
		err := decodeJSONBody(w, r, &got, 131072)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidRequest.Error())
	})

	t.Run("empty body", func(t *testing.T) {
		w := setup(t)
		r := httptest.NewRequest(http.MethodGet, "/1", nil)
		r.Header.Set("Content-Type", "application/json")
		err := decodeJSONBody(w, r, nil, 131072)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidRequest.Error())
	})

	t.Run("large body", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo string `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`{"foo":"very long body"}`))
		r.Header.Set("Content-Type", "application/json")
		err := decodeJSONBody(w, r, &got, 1)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidRequest.Error())
	})

	t.Run("multiple json objects", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo string `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`{"foo":"very long body"}{"bar":10}`))
		r.Header.Set("Content-Type", "application/json")
		err := decodeJSONBody(w, r, &got, 131072)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidRequest.Error())
	})

	t.Run("successful decode", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo string   `json:"foo"`
			Bar int      `json:"bar"`
			Baz []string `json:"baz"`
		}
		var got body
		b := `{
  "foo": "value",
  "bar": 10,
  "baz": [
    "value1",
    "value2",
    "value3"
  ]
}`
		want := body{
			Foo: "value",
			Bar: 10,
			Baz: []string{"value1", "value2", "value3"},
		}
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(b))
		r.Header.Set("Content-Type", "application/json")
		err := decodeJSONBody(w, r, &got, 131072)

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})
}

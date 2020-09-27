/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package httputil

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"adeia"
	"adeia/pkg/errs"

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
		LogWriteErr(m, errors.New("test error"))
		want := fmt.Errorf("failed to write response: %v", errors.New("test error"))
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
		{adeia.ErrInvalidJSON, true, "req malformed error"},
		{adeia.ErrInvalidRequestBody, true, "req malformed error"},
		{adeia.ErrValidationFailed, true, "req malformed error"},
		{adeia.ErrRequestBodyTooLarge, true, "req malformed error"},
		{adeia.ErrUnknownField, true, "req malformed error"},
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

func TestRespond(t *testing.T) {
	t.Run("json marshal error: UnsupportedTypeError", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		payload := make(chan int)
		err := Respond(w, 0, payload)
		assert.Error(t, err)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
		assert.Equal(t, strings.TrimSuffix(w.Body.String(), "\n"), http.StatusText(http.StatusInternalServerError))
	})

	t.Run("json marshal error: UnsupportedValueError", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		payload := math.Inf(1)
		err := Respond(w, 0, payload)

		assert.Error(t, err)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
		assert.Equal(t, strings.TrimSuffix(w.Body.String(), "\n"), http.StatusText(http.StatusInternalServerError))
	})

	t.Run("error writing response", func(t *testing.T) {
		t.Parallel()

		w := &errResponseWriter{httptest.NewRecorder()}
		payload := struct {
			Foo string
			Bar string
		}{"foo", "bar"}
		err := Respond(w, 100, payload)

		assert.Error(t, err)
		assert.EqualError(t, err, errFailedToWriteResponse.Error())
		assert.Equal(t, w.Code, 100)
		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
	})

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()

		want := `{"foo":"foo","bar":"bar"}`
		w := httptest.NewRecorder()
		payload := struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		}{"foo", "bar"}
		err := Respond(w, http.StatusOK, payload)
		resp := w.Result()
		got, _ := ioutil.ReadAll(resp.Body)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
		assert.Equal(t, want, string(got))
	})

	t.Run("respond with data", func(t *testing.T) {
		t.Parallel()

		want := `{"data":{"foo":"foo","bar":"bar"}}`
		w := httptest.NewRecorder()
		payload := struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		}{"foo", "bar"}
		err := RespondWithData(w, http.StatusOK, payload)
		resp := w.Result()
		got, _ := ioutil.ReadAll(resp.Body)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
		assert.Equal(t, want, string(got))
	})

	t.Run("respond with error", func(t *testing.T) {
		t.Parallel()

		want := `{"error":{"code":"TEST_ERROR_CODE","message":"test error"}}`
		w := httptest.NewRecorder()
		payload := errs.ResponseError{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  "TEST_ERROR_CODE",
			Message:    "test error",
		}
		err := RespondWithErr(w, payload)
		resp := w.Result()
		got, _ := ioutil.ReadAll(resp.Body)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
		assert.Equal(t, want, string(got))
	})
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
		err := DecodeJSONBody(w, r, nil)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidRequestBody.Error())
	})

	t.Run("syntax error in json", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo int `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`foo`))
		err := DecodeJSONBody(w, r, &got)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidJSON.Error())
	})

	t.Run("json validation failed", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo int `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`{"foo":["bar","baz"]}`))
		err := DecodeJSONBody(w, r, &got)

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
		err := DecodeJSONBody(w, r, &got)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrUnknownField.Msg("Unknown field: bar").Error())
	})

	t.Run("empty body", func(t *testing.T) {
		w := setup(t)
		r := httptest.NewRequest(http.MethodGet, "/1", nil)
		err := DecodeJSONBody(w, r, nil)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidJSON.Error())
	})

	t.Run("large body", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo string `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`{"foo":"very long body"}`))
		err := decodeJSONBody(w, r, &got, 1)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrRequestBodyTooLarge.Error())
	})

	t.Run("multiple json objects", func(t *testing.T) {
		w := setup(t)
		type body struct {
			Foo string `json:"foo"`
		}
		var got body
		r := httptest.NewRequest(http.MethodGet, "/1", strings.NewReader(`{"foo":"very long body"}{"bar":10}`))
		err := DecodeJSONBody(w, r, &got)

		assert.Error(t, err)
		assert.EqualError(t, err, adeia.ErrInvalidJSON.Error())
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
		err := DecodeJSONBody(w, r, &got)

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})
}

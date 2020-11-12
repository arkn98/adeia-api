/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package httputil

import (
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"adeia/pkg/http/response"

	"github.com/stretchr/testify/assert"
)

func TestRespond(t *testing.T) {
	t.Run("json marshal error: UnsupportedTypeError", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		payload := make(chan int)
		err := Respond(w, response.NewResponse(response.WithStatusCode(0), response.WithData(payload)))
		assert.Error(t, err)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
		assert.Equal(t, strings.TrimSuffix(w.Body.String(), "\n"), http.StatusText(http.StatusInternalServerError))
	})

	t.Run("json marshal error: UnsupportedValueError", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		payload := math.Inf(1)
		err := Respond(w, response.NewResponse(response.WithStatusCode(0), response.WithData(payload)))

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
		err := Respond(w, response.NewResponse(response.WithStatusCode(100), response.WithData(payload)))

		assert.Error(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "failed to write response"))
		assert.Equal(t, w.Code, 100)
		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
	})

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()

		want := `{"data":{"foo":"foo","bar":"bar"}}`
		w := httptest.NewRecorder()
		payload := struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		}{"foo", "bar"}
		err := Respond(w, response.NewResponse(response.WithStatusCode(http.StatusOK), response.WithData(payload)))
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
		err := Respond(w, response.NewResponse(response.WithStatusCode(http.StatusOK), response.WithData(payload)))
		resp := w.Result()
		got, _ := ioutil.ReadAll(resp.Body)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
		assert.Equal(t, want, string(got))
	})

	t.Run("respond with error", func(t *testing.T) {
		t.Parallel()

		want := `{"error":{"error_type":"TEST_ERROR","code":"TEST_ERROR_CODE","message":"test error"}}`
		w := httptest.NewRecorder()
		payload := &response.Error{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorType:      "TEST_ERROR",
			ErrorCode:      "TEST_ERROR_CODE",
			Message:        "test error",
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

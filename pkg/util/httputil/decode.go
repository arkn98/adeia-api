/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"adeia"
	"adeia/pkg/constants"
	"adeia/pkg/http/response"

	"github.com/golang/gddo/httputil/header"
)

const (
	statusRequestEntityTooLargeMessage = "http: request body too large"
	statusJSONUnknownField             = "json: unknown field"
)

// decodeJSONBody decodes a JSON http.Request.Body into the provided interface.
// `dest` must be a pointer so that the decoded value can be used. An error is
// returned when parsing cannot happen.
//
// Adapted from Alex Edwards's blog, released under the MIT license.
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dest interface{}, maxBodySize int64) error {
	value, params := header.ParseValueAndParams(r.Header, "Content-Type")
	if value != "application/json" || len(params) != 0 {
		// reject all other content-types or when other params are present
		return adeia.ErrUnsupportedMediaType
	}

	// set max body size; err is returned when body exceeds this size
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case
			errors.As(err, &syntaxError),                           // bad JSON
			errors.Is(err, io.ErrUnexpectedEOF),                    // unexpected EOF
			errors.Is(err, io.EOF),                                 // request body is empty
			strings.HasPrefix(err.Error(), statusJSONUnknownField), // unknown fields are present
			err.Error() == statusRequestEntityTooLargeMessage:      // request body is too large
			return adeia.ErrInvalidRequest

		case errors.As(err, &unmarshalTypeError): // invalid value for field
			return adeia.ErrValidationFailed.
				AddValidationErr(
					unmarshalTypeError.Field,
					fmt.Sprintf("Please enter a valid value for %v", unmarshalTypeError.Field),
				)

		default: // unmarshal-able target
			return err
		}
	}

	if dec.More() {
		// body contains multiple JSON objects
		return adeia.ErrInvalidRequest
	}

	return nil
}

func isReqMalformedErr(err error) bool {
	switch err.Error() {
	case
		adeia.ErrInvalidRequest.Error(),
		adeia.ErrUnsupportedMediaType.Error(),
		adeia.ErrValidationFailed.Error():
		return true
	}

	return false
}

// Decode is a wrapper around decodeJSONBody that decodes the request body into
// a destination interface. And if decodeJSONBody returns an error, an appropriate
// error response is sent back as response. So, the controller calling this method
// need not write a response.
func Decode(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	if err := decodeJSONBody(w, r, dest, constants.MaxReqBodySize); err != nil {
		if isReqMalformedErr(err) {
			// we ignore write errors here, so a response is not guaranteed
			_ = RespondWithErr(w, err.(*response.Error))
			return fmt.Errorf("malformed request body: %v", err)
		}

		// some other error
		_ = RespondWithErr(w, adeia.ErrParseReqBodyFailed)
		return fmt.Errorf("cannot parse request body: %v", err)
	}

	return nil
}

// LogWarner is an interface for a logger that can Warnf(). This is specifically used
// when logging errors that occur when writing a HTTP response.
type LogWarner interface {
	Warnf(template string, args ...interface{})
}

// LogWriteErr logs the errors that occur when http.ResponseWriter.Write() returns
// an error. When Write() returns an error, the least we can do is to log it, so that
// it can be looked up in the future.
func LogWriteErr(log LogWarner, err error) {
	if err != nil {
		// we cannot do anything if we can't send a response, so we just log
		log.Warnf("failed to write response: %v", err)
	}
}

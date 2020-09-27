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
	"strconv"
	"strings"

	"adeia"
	"adeia/pkg/errs"
	"adeia/pkg/util/constants"

	"github.com/golang/gddo/httputil/header"
)

var errFailedToWriteResponse = errors.New("failed to write response")

const (
	statusRequestEntityTooLargeMessage = "http: request body too large"
	statusJSONUnknownField             = "json: unknown field "
)

// DecodeJSONBody decodes a JSON http.Request.Body into the provided interface, dest,
// (usually a struct). dest must be a pointer so that the filled-in values can be used.
// An error is returned when the parsing cannot happen.
//
// Adapted from Alex Edwards's blog, released under the MIT license.
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	return decodeJSONBody(w, r, dest, constants.MaxReqBodySize)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dest interface{}, maxBodySize int64) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			// request body is not JSON
			return adeia.ErrInvalidRequestBody
		}
	}

	// set max body size; err is returned when body exceeds this size
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	// TODO: decide if we need to disallow unknown fields
	dec.DisallowUnknownFields()

	err := dec.Decode(&dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case
			errors.As(err, &syntaxError),
			errors.Is(err, io.ErrUnexpectedEOF):
			// badly formed JSON
			return adeia.ErrInvalidJSON

		case errors.As(err, &unmarshalTypeError):
			// invalid value for field
			return adeia.ErrValidationFailed.
				AddValidationErr(
					unmarshalTypeError.Field,
					fmt.Sprintf("Please enter a valid %v", unmarshalTypeError.Field),
				)

		// There is an open issue regarding turning this into a sentinel error
		// at https://github.com/golang/go/issues/29035.
		case strings.HasPrefix(err.Error(), statusJSONUnknownField):
			// unknown field is present
			field, _ := strconv.Unquote(strings.TrimPrefix(err.Error(), "json: unknown field "))
			return adeia.ErrUnknownField.Msgf("Unknown field: %v", field)

		case errors.Is(err, io.EOF):
			// request body is empty
			return adeia.ErrInvalidJSON

		// There is an open issue regarding turning this into a sentinel error
		// at https://github.com/golang/go/issues/30715.
		case err.Error() == statusRequestEntityTooLargeMessage:
			// request body is too large
			return adeia.ErrRequestBodyTooLarge

		default:
			return err
		}
	}

	// we see if any remaining body is pending
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		// contains multiple JSON objects
		return adeia.ErrInvalidJSON
	}

	return nil
}

func isReqMalformedErr(err error) bool {
	switch err.Error() {
	case
		adeia.ErrInvalidJSON.Error(),
		adeia.ErrInvalidRequestBody.Error(),
		adeia.ErrValidationFailed.Error(),
		adeia.ErrRequestBodyTooLarge.Error(),
		adeia.ErrUnknownField.Error():
		return true
	}

	return false
}

// Decode is a wrapper around DecodeJSONBody that decodes the request body into
// a destination interface. And if DecodeJSONBody returns an error, an appropriate
// error response is sent back. So the controller calling this method need not
// write a response. But, the errors returned when writing the error response are
// ignored; so in that case, a response is not 100% guaranteed.
func Decode(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	if err := DecodeJSONBody(w, r, dest); err != nil {
		if isReqMalformedErr(err) {
			_ = RespondWithErr(w, err.(errs.ResponseError))
			return fmt.Errorf("malformed request body: %v", err)
		}

		// some other error
		_ = RespondWithErr(w, adeia.ErrParseReqBodyFailed)
		return fmt.Errorf("cannot parse request body: %v", err)
	}

	return nil
}

// Respond is a util that writes a HTTP statusCode and payload, as a HTTP JSON response.
func Respond(w http.ResponseWriter, statusCode int, payload interface{}) error {
	resp, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("failed to marshal payload to JSON: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(resp)
	if err != nil {
		// we use this to log
		return errFailedToWriteResponse
	}
	return nil
}

type errorResponse struct {
	Error errs.ResponseError `json:"error"`
}

type dataResponse struct {
	Data interface{} `json:"data"`
}

// RespondWithErr is a wrapper around Respond that structures the response to be
// a error response, with fields like error code, validation errors, message, etc.
func RespondWithErr(w http.ResponseWriter, err errs.ResponseError) error {
	payload := &errorResponse{err}
	return Respond(w, err.StatusCode, payload)
}

// RespondWithData is a wrapper around Respond that structures the response to be
// a data response.
func RespondWithData(w http.ResponseWriter, statusCode int, data interface{}) error {
	payload := &dataResponse{data}
	return Respond(w, statusCode, payload)
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

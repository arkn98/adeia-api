/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package errs

import "fmt"

// ResponseError represents an error that is sent as a response to the client.
// A ResponseError is usually a more generic version of an error, that can be
// safely sent to the client without leaking any internal details.
type ResponseError struct {
	// StatusCode is the HTTP response code that should be returned for the error.
	StatusCode int `json:"-"`

	// ErrorCode is a constant string representing the type of error. It must be
	// unique and must not change, as it is used in comparisons.
	ErrorCode string `json:"code"`

	// Message is an optional short message for the user, describing the details
	// of the error.
	Message string `json:"message,omitempty"`

	// ValidationErrors is a map of validation errors, with keys as fields and
	// corresponding error messages of the fields as values.
	ValidationErrors map[string]string `json:"validation_errors,omitempty"`
}

// Msg sets ResponseError's Message.
func (re ResponseError) Msg(m string) ResponseError {
	re.Message = m
	return re
}

// Msgf sets ResponseError's Message as a formatted string.
func (re ResponseError) Msgf(format string, a ...interface{}) ResponseError {
	re.Message = fmt.Sprintf(format, a...)
	return re
}

// Error returns the ErrorCode of the ResponseError.
func (re ResponseError) Error() string {
	return re.ErrorCode
}

// AddValidationErr adds a new validation error to the ValidationErrors map.
func (re ResponseError) AddValidationErr(f, m string) ResponseError {
	if re.ValidationErrors == nil {
		re.ValidationErrors = make(map[string]string)
	}
	re.ValidationErrors[f] = m
	return re
}

// ValidationErr sets the entire ValidationErrors map to the provided map, which
// is useful when processing fields using third-party validation libraries.
func (re ResponseError) ValidationErr(m map[string]string) ResponseError {
	re.ValidationErrors = m
	return re
}

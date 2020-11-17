/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package response

import "fmt"

// Error represents an error that is sent as a response to the client.
// A Error is usually a more generic version of an error, that can be
// safely sent to the client without revealing any internal details.
type Error struct {
	// HTTPStatusCode is the HTTP response code that should be returned for the error.
	HTTPStatusCode int `json:"-"`

	// Type represents the type of the error returned.
	ErrorType string `json:"error_type"`

	// ErrorCode is a constant string representing the type of error. It must be
	// unique and must not change, as it is used in comparisons.
	ErrorCode string `json:"code"`

	// Message is a short message describing the details of the error.
	Message string `json:"message"`

	// ValidationErrors is a map of validation errors, with keys as fields and
	// corresponding error messages of the fields as values.
	ValidationErrors map[string]string `json:"validation_errors,omitempty"`
}

// NewError creates a new *Error.
func NewError() *Error {
	return &Error{}
}

// StatusCode sets the Error's HTTPStatusCode.
func (e *Error) StatusCode(c int) *Error {
	e.HTTPStatusCode = c
	return e
}

// Msg sets Error's Message.
func (e *Error) Msg(m string) *Error {
	e.Message = m
	return e
}

// Msgf sets Error's Message as a formatted string.
func (e *Error) Msgf(format string, args ...interface{}) *Error {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// Type sets the Error's ErrorType.
func (e *Error) Type(t string) *Error {
	e.ErrorType = t
	return e
}

// Code set's the Errors' ErrorCode.
func (e *Error) Code(c string) *Error {
	e.ErrorCode = c
	return e
}

// Error returns the ErrorCode.
func (e *Error) Error() string {
	return e.ErrorCode
}

// AddValidationErr adds a new validation error to the ValidationErrors map.
func (e *Error) AddValidationErr(f, m string) *Error {
	if e.ValidationErrors == nil {
		e.ValidationErrors = make(map[string]string)
	}
	e.ValidationErrors[f] = m
	return e
}

// ValidationErr sets the entire ValidationErrors map to the provided map, which
// is useful when processing fields using third-party validation libraries.
func (e *Error) ValidationErr(m map[string]string) *Error {
	e.ValidationErrors = m
	return e
}

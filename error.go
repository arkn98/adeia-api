/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package adeia

import (
	"net/http"

	"adeia/pkg/http/response"
)

var (
	// ErrUnsupportedMediaType is the error returned when the request `Content-Type`
	// is not supported by the server.
	ErrUnsupportedMediaType = response.NewError().
				StatusCode(http.StatusUnsupportedMediaType).
				Code("UNSUPPORTED_MEDIA_TYPE").
				Msg("Content-Type must be application/json without any parameters")

	// ErrInvalidRequest is the error returned when the request is bad.
	ErrInvalidRequest = response.NewError().
				StatusCode(http.StatusBadRequest).
				Code("INVALID_REQUEST").
				Msg("Request body is malformed")

	// ErrValidationFailed is the error returned when some of the fields do not conform to
	// the validation rules. Add a list of ValidationErrors to this to specify the fields
	// that are failing.
	ErrValidationFailed = response.NewError().
				StatusCode(http.StatusBadRequest).
				Code("VALIDATION_FAILED").
				Msg("Validation failed for some fields")

	// ErrParseReqBodyFailed is the error returned when the request body is okay,
	// but something else happened and we cannot parse it.
	ErrParseReqBodyFailed = response.NewError().
				StatusCode(http.StatusInternalServerError).
				Code("PARSE_REQUEST_BODY_FAILED").
				Msg("An error occurred while parsing the request body")

	// ErrAPIError is the error returned when an error occurs on the API side. It
	// is made as generic as possible, so that internal details are not revealed outside.
	ErrAPIError = response.NewError().
			StatusCode(http.StatusInternalServerError).
			Code("API_ERROR").
			Msg("An error occurred while parsing the request body")

	// ErrResourceAlreadyExists is the error returned when a resource already
	// exists with the specified fields.
	ErrResourceAlreadyExists = response.NewError().
					StatusCode(http.StatusBadRequest).
					Code("RESOURCE_ALREADY_EXISTS").
					Msg("A resource already exists with the specified fields")
)

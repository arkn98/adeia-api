package adeia

import (
	"net/http"

	"adeia/pkg/errs"
)

var (
	// ErrInvalidRequestBody is the error returned when the request body is not JSON.
	ErrInvalidRequestBody = errs.ResponseError{
		StatusCode: http.StatusUnsupportedMediaType,
		ErrorCode:  "INVALID_REQUEST_BODY",
		Message:    "Request body must be JSON",
	}

	// ErrInvalidJSON is the error returned when the request body contains invalid JSON.
	ErrInvalidJSON = errs.ResponseError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "INVALID_JSON",
		Message:    "An error occurred while parsing the request body",
	}

	// ErrValidationFailed is the error returned when some of the fields do not conform to
	// the validation rules. Add a list of ValidationErrors to this to specify the fields
	// that are failing.
	ErrValidationFailed = errs.ResponseError{
		StatusCode: http.StatusUnprocessableEntity,
		ErrorCode:  "VALIDATION_FAILED",
		Message:    "Validation failed for some fields",
	}

	// ErrRequestBodyTooLarge is the error returned when the request body is too large.
	ErrRequestBodyTooLarge = errs.ResponseError{
		StatusCode: http.StatusRequestEntityTooLarge,
		ErrorCode:  "REQUEST_BODY_TOO_LARGE",
		Message:    "Request body is too large",
	}

	// ErrUnknownField is the error returned when the request body contains an unknown field.
	ErrUnknownField = errs.ResponseError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "UNKNOWN_FIELD",
		Message:    "An unknown field is present in the request body",
	}

	// ErrDatabaseError is the error returned when a database error occurs.
	// It is made as generic as possible, so that internals are not revealed outside.
	ErrDatabaseError = errs.ResponseError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  "DATABASE_ERROR",
	}

	// ErrResourceAlreadyExists is the error returned when a resource already
	// exists with the specified fields.
	ErrResourceAlreadyExists = errs.ResponseError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "RESOURCE_ALREADY_EXISTS",
		Message:    "A resource already exists with the specified fields",
	}

	// ErrParseReqBodyFailed is the error returned when the request body is okay,
	// but something else happened and we cannot parse it.
	ErrParseReqBodyFailed = errs.ResponseError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  "PARSE_REQUEST_BODY_FAILED",
		Message:    "An error occurred while parsing the request body",
	}
)

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package response

// Response represents the generalized response that will be sent to the client.
type Response struct {
	// StatusCode is the HTTP response code that should be sent for the response.
	StatusCode int         `json:"-"`
	ID         string      `json:"id,omitempty"`
	ObjectType string      `json:"object,omitempty"`
	HasMore    bool        `json:"has_more,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      *Error      `json:"error,omitempty"`
}

// Opt represents the optional function to modify the Response.
type Opt func(r *Response)

// WithStatusCode sets the Response's StatusCode field.
func WithStatusCode(code int) Opt {
	return func(r *Response) {
		r.StatusCode = code
	}
}

// WithID sets the Response's ID field.
func WithID(id string) Opt {
	return func(r *Response) {
		r.ID = id
	}
}

// WithObjectType sets the Response's ObjectType field.
func WithObjectType(t string) Opt {
	return func(r *Response) {
		r.ObjectType = t
	}
}

// WithHasMore sets the Response's HasMore field.
func WithHasMore(b bool) Opt {
	return func(r *Response) {
		r.HasMore = b
		r.ObjectType = "list"
	}
}

// HasMore wraps WithHasMore.
func HasMore() Opt {
	return WithHasMore(true)
}

// WithData sets the Response's Data field.
func WithData(d interface{}) Opt {
	return func(r *Response) {
		r.Data = d
	}
}

// WithError sets the Response's error field.
func WithError(e *Error) Opt {
	return func(r *Response) {
		r.Error = e
		r.StatusCode = e.HTTPStatusCode
	}
}

// NewResponse creates a new *Response that can be sent to the client.
func NewResponse(opts ...Opt) *Response {
	r := &Response{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"adeia/pkg/http/response"
)

// RespondWithErr is a wrapper around Respond for sending errors as a response.
func RespondWithErr(w http.ResponseWriter, e *response.Error) error {
	return Respond(w, response.NewResponse(response.WithError(e)))
}

// Respond is a util that writes a HTTP statusCode and payload, as a HTTP JSON response.
func Respond(w http.ResponseWriter, r *response.Response) error {
	resp, err := json.Marshal(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("failed to marshal payload to JSON: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	_, err = w.Write(resp)
	if err != nil {
		return errors.New("failed to write response")
	}
	return nil
}

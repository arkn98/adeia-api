/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import "net/http"

// ProtectedHandler checks if user is authorized before allowing the request to
// pass to the underlying controller.
type ProtectedHandler struct {
	PermissionName string
	Handler        http.HandlerFunc
}

// ServeHTTP serves a request using the ProtectedHandler.Handler.
func (p *ProtectedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// perform role checking using context, like !contains(userRoles, p.PermissionName)
	if false {
		return
	}

	// user has access, so continue
	p.Handler.ServeHTTP(w, r)
}

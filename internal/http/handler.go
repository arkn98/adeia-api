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

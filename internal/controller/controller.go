package controller

import (
	"net/http"

	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	holidayService "adeia-api/internal/service/holiday"
	roleService "adeia-api/internal/service/role"
	sessionService "adeia-api/internal/service/session"
	userService "adeia-api/internal/service/user"
	"adeia-api/internal/util"
	"adeia-api/internal/util/log"
	"adeia-api/internal/util/mail"
)

// ProtectedHandler checks if user is authorized before allowing the request to
// pass to the underlying controller.
type ProtectedHandler struct {
	PermissionName string
	Handler        http.HandlerFunc
}

// ServeHTTP performs the authorization and only allows the request to pass when
// the user is authorized.
func (p *ProtectedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// perform role checking using context, like !contains(userRoles, p.PermissionName)
	if false {
		// user doesn't have access
		log.Debug("user not authorized to perform action")
		util.RespondWithError(w, util.ErrUnauthorized)
		return
	}

	// user has access, so continue
	p.Handler.ServeHTTP(w, r)
}

var (
	holidaySvc holidayService.Service
	roleSvc    roleService.Service
	sessionSvc sessionService.Service
	usrSvc     userService.Service
)

// Init initializes all services that are used in the controllers.
func Init(d db.DB, c cache.Cache, m mail.Mailer) {
	holidaySvc = holidayService.New(d, c)
	roleSvc = roleService.New(d, c)
	sessionSvc = sessionService.New(d)
	usrSvc = userService.New(d, c, m)
}

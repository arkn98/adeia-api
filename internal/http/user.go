/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"fmt"
	"net/http"

	"adeia"
	"adeia/pkg/errs"
	"adeia/pkg/log"
	"adeia/pkg/util/constants"
	"adeia/pkg/util/httputil"

	"github.com/go-chi/chi"
)

// UserController represents the User controller.
type UserController struct {
	handler     chi.Router
	log         log.Logger
	pattern     string
	userService adeia.UserService
}

// Handler returns the UserController's handler.
func (uc *UserController) Handler() http.Handler {
	return uc.handler
}

// Pattern returns the UserController's pattern.
func (uc *UserController) Pattern() string {
	return uc.pattern
}

// NewUserController creates a new UserController.
func NewUserController(log log.Logger, us adeia.UserService) *UserController {
	uc := &UserController{
		log:         log,
		pattern:     "/users",
		userService: us,
	}
	uc.BindRoutes()
	return uc
}

// BindRoutes binds all user-routes to the UserController's handler.
func (uc *UserController) BindRoutes() {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/", uc.CreateUser())
	//r.Method(http.MethodGet, "/", uc.CheckContext())

	uc.handler = r
}

//func (uc *UserController) CheckContext() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		fmt.Println("request started")
//		defer fmt.Println("request ended")
//
//		select {
//		case <-time.After(10 * time.Second):
//			fmt.Println("hello")
//		case <-ctx.Done():
//			err := ctx.Err()
//			fmt.Println(err)
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//		}
//	}
//}

// CreateUser creates a new User if it doesn't exist already.
func (uc *UserController) CreateUser() *ProtectedHandler {
	type request struct {
		Name        string `json:"name"`
		EmployeeID  string `json:"employee_id,omitempty"`
		Email       string `json:"email"`
		Designation string `json:"designation"`
	}

	return &ProtectedHandler{
		PermissionName: "CREATE_USERS",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var body request
			if err := httputil.Decode(w, r, &body); err != nil {
				uc.log.Debug(err)
				return
			}

			// validate request

			user, err := uc.userService.CreateUser(
				r.Context(),
				body.Name,
				body.Email,
				body.EmployeeID,
				body.Designation,
			)
			if err != nil {
				httputil.LogWriteErr(uc.log, httputil.RespondWithErr(w, err.(errs.ResponseError)))
				return
			}

			w.Header().Set("Location", fmt.Sprintf(
				"/%s/%s/%s",
				constants.APIVersion,
				uc.pattern,
				user.EmployeeID,
			))
			httputil.LogWriteErr(uc.log, httputil.RespondWithData(w, http.StatusCreated, user))
		},
	}
}

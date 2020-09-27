/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package service

import (
	"context"

	"adeia"
	"adeia/pkg/log"
	"adeia/pkg/util/crypto"
)

// UserService represents the User service.
type UserService struct {
	log  log.Logger
	repo adeia.UserRepo
}

// NewUserService creates a new *UserService.
func NewUserService(log log.Logger, repo adeia.UserRepo) *UserService {
	return &UserService{log, repo}
}

// CreateUser creates a new user if does not exist.
func (us *UserService) CreateUser(ctx context.Context, name, email, empID, designation string) (*adeia.User, error) {
	if u, err := us.repo.GetByEmail(ctx, email); err != nil {
		us.log.Errorf("cannot fetch user by email: %v", err)
		return nil, adeia.ErrDatabaseError
	} else if u != nil {
		us.log.Debug("user already exists with the provided email " + email)
		return nil, adeia.ErrResourceAlreadyExists
	}

	if empID == "" {
		empID = crypto.NewEmpID()
	}
	user := adeia.NewUser(
		adeia.WithName(name),
		adeia.WithEmail(email),
		adeia.WithDesignation(designation),
		adeia.WithEmpID(empID),
	)

	if _, err := us.repo.Insert(ctx, user); err != nil {
		us.log.Warnf("cannot create new user: %v", err)
		return nil, adeia.ErrDatabaseError
	}
	return user, nil
}

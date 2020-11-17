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
func NewUserService(l log.Logger, r adeia.UserRepo) *UserService {
	return &UserService{l, r}
}

// CreateUser creates a new user if it does not exist.
func (s *UserService) CreateUser(ctx context.Context, name, email, empID, designation string) (*adeia.User, error) {
	if u, err := s.repo.GetByEmail(ctx, email); err != nil {
		s.log.Errorf("cannot fetch user by email: %v", err)
		return nil, adeia.ErrAPIError
	} else if u != nil {
		// TODO: do not reveal that user account already exists
		s.log.Debugf("user already exists with the provided email: %v ", email)
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

	if _, err := s.repo.Insert(ctx, user); err != nil {
		s.log.Warnf("cannot create new user: %v", err)
		return nil, adeia.ErrAPIError
	}
	return user, nil
}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package repo

import (
	"context"

	"adeia"
	"adeia/internal/store"
)

const (
	queryByEmail = "SELECT * FROM users WHERE email=$1"
	//queryByEmpID            = "SELECT * FROM users WHERE employee_id=$1"
	//queryByID               = "SELECT * FROM users WHERE id=$1"
	queryInsert = "INSERT INTO users (employee_id, name, email, password, designation, is_activated) VALUES " +
		"(:employee_id, :name, :email, :password, :designation, :is_activated) RETURNING id"
	//queryUpdatePwdAndIsActivated = "UPDATE users SET password=:password, is_activated=:is_activated " +
	//	"WHERE id=:id"
)

// UserRepo represents the User repository.
type UserRepo struct {
	db store.DB
}

// NewUserRepo creates a new *UserRepo.
func NewUserRepo(d store.DB) *UserRepo {
	return &UserRepo{d}
}

//
//func (ur *UserRepo) DeleteByEmpID(empID string) (rowsAffected int64, err error) {
//	return ur.db.Delete(queryDeleteByEmpID, time.Now().UTC(), empID)
//}
//
//func (ur *UserRepo) GetByEmail(email string) (*adeia.User, error) {
//	return ur.get(queryByEmail, email)
//}
//
//func (ur *UserRepo) GetByEmailInclDeleted(email string) (*adeia.User, error) {
//	return ur.get(queryByEmailInclDeleted, email)
//}
//
//func (ur *UserRepo) GetByEmpID(empID string) (*adeia.User, error) {
//	return ur.get(queryByEmpID, empID)
//}
//
//func (ur *UserRepo) GetByID(id int) (*adeia.User, error) {
//	return ur.get(queryByID, id)
//}

// Insert inserts a new User and returns the lastInsertID.
func (ur *UserRepo) Insert(ctx context.Context, u *adeia.User) (lastInsertID int, err error) {
	return ur.db.InsertNamed(ctx, queryInsert, u)
}

// GetByEmail returns a User using the provided email address.
func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (*adeia.User, error) {
	return ur.get(ctx, queryByEmail, email)
}

//func (ur *UserRepo) UpdatePasswordAndIsActivated(u *adeia.User, password string, isActivated bool) error {
//	u.Password = password
//	u.IsActivated = isActivated
//	if _, err := ur.db.UpdateNamed(queryUpdatePwdAndIsActivated, u); err != nil {
//		return err
//	}
//	return nil
//}

func (ur *UserRepo) get(ctx context.Context, query string, args ...interface{}) (*adeia.User, error) {
	u := adeia.User{}
	if ok, err := ur.db.GetOne(ctx, &u, query, args...); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &u, nil
}

func (ur *UserRepo) getMany(ctx context.Context, query string, args ...interface{}) ([]*adeia.User, error) {
	var u []*adeia.User
	if err := ur.db.GetMany(ctx, &u, query, args...); err != nil {
		return nil, err
	} else if len(u) == 0 {
		return nil, nil
	}
	return u, nil
}

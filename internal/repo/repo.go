package repo

import (
	"adeia-api/internal/model"
)

// UserRepo is an interface that represents the list of functions that need to be
// implemented for the User model, by the repo.
type UserRepo interface {
	GetByEmail(email string) (*model.User, error)
	GetByEmpID(empID string) (*model.User, error)
	GetByID(id int) (*model.User, error)
	GetByEmpIDAndEmail(empID, email string) (*model.User, error)
	Insert(u *model.User) (int, error)
	UpdatePasswordAndIsActivated(u *model.User, password string, isActivated bool) error
}

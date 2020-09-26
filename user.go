package adeia

import "context"

// User represents the User model.
type User struct {
	// ID is a surrogate primary key that is auto-incremented at the database.
	// This has no meaning outside the database, except that it only identifies an account.
	// This will never change for an account and all foreign keys must use this field.
	// This must not be exposed to the outside.
	ID int `db:"id" json:"-"`

	// EmployeeID is a natural key used extensively throughout the system (in URIs, etc.).
	// It must be unique, short and user-rememberable (preferably 6-8 chars long).
	// It is case-insensitive (internally managed by Postgres as `citext` (case-insensitive text)).
	// https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#user-ids
	EmployeeID string `db:"employee_id" json:"employee_id"`

	// Name represents the name of the User.
	Name string `db:"name" json:"name"`

	// Email represents the email of the User.
	Email string `db:"email" json:"email"`

	// Password represents the hashed & salted password of the User.
	Password string `db:"password" json:"-"`

	// Designation represents the designation of the User.
	Designation string `db:"designation" json:"designation"`

	// IsActivated represents whether the User account is activated or not.
	IsActivated bool `db:"is_activated" json:"is_activated"`
}

// UserRepo is the interface for all the repository functions on the User model.
type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Insert(ctx context.Context, u *User) (lastInsertID int, err error)
}

// UserService is the interface for all the business rules on the User model.
type UserService interface {
	CreateUser(ctx context.Context, name, email, empID, designation string) (*User, error)
}

// UserOpt represents the optional function to modify the User.
type UserOpt func(u *User)

// WithName is an UserOpt to set the Name of the User.
func WithName(n string) UserOpt {
	return func(u *User) {
		u.Name = n
	}
}

// WithEmpID is an UserOpt to set the EmployeeID of the User.
func WithEmpID(e string) UserOpt {
	return func(u *User) {
		u.EmployeeID = e
	}
}

// WithEmail is an UserOpt to set the Email of the User.
func WithEmail(e string) UserOpt {
	return func(u *User) {
		u.Email = e
	}
}

// WithDesignation is an UserOpt to set the Designation of the User.
func WithDesignation(d string) UserOpt {
	return func(u *User) {
		u.Designation = d
	}
}

// WithActivation is an UserOpt to mark the User as activated.
func WithActivation() UserOpt {
	return func(u *User) {
		u.IsActivated = true
	}
}

// WithPassword is an UserOpt to set the password of the User.
func WithPassword(p string) UserOpt {
	return func(u *User) {
		u.Password = p
	}
}

// NewUser creates a new User, with a set of defaults. The defaults can be overridden
// using the various UserOpts.
func NewUser(opts ...UserOpt) *User {
	u := &User{
		IsActivated: false,
		Password:    "",
	}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

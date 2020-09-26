package adeia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithName(t *testing.T) {
	t.Parallel()
	want := "foobar"
	opt := WithName(want)
	u := &User{}
	opt(u)
	assert.Equal(t, want, u.Name)
}

func TestWithEmpID(t *testing.T) {
	t.Parallel()
	want := "foobar"
	opt := WithEmpID(want)
	u := &User{}
	opt(u)
	assert.Equal(t, want, u.EmployeeID)
}

func TestWithEmail(t *testing.T) {
	t.Parallel()
	want := "foobar"
	opt := WithEmail(want)
	u := &User{}
	opt(u)
	assert.Equal(t, want, u.Email)
}

func TestWithDesignation(t *testing.T) {
	t.Parallel()
	want := "foobar"
	opt := WithDesignation(want)
	u := &User{}
	opt(u)
	assert.Equal(t, want, u.Designation)
}

func TestWithIsActivated(t *testing.T) {
	t.Parallel()
	opt := WithActivation()
	u := &User{}
	opt(u)
	assert.True(t, u.IsActivated)
}

func TestWithPassword(t *testing.T) {
	t.Parallel()
	want := "foobar"
	opt := WithPassword(want)
	u := &User{}
	opt(u)
	assert.Equal(t, want, u.Password)
}

func TestNewUser(t *testing.T) {
	t.Run("default user", func(t *testing.T) {
		t.Parallel()
		want := &User{
			IsActivated: false,
			Password:    "",
		}
		got := NewUser()
		assert.Equal(t, want, got, "default user should not be activated and have empty password")
	})

	t.Run("with opts", func(t *testing.T) {
		t.Parallel()
		want := &User{
			IsActivated: true,
			Password:    "",
			Email:       "foobar",
		}
		got := NewUser(
			WithActivation(),
			WithEmail("foobar"),
		)
		assert.Equal(t, want, got)
	})
}

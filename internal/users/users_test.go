package users

import (
	"testing"

	"github.com/mithcs/probox-api/pkg/am"
)

func TestValidateUsername(t *testing.T) {
	t.Run("valid username", func(t *testing.T) {
		username := "valid_username"
		err := validateUsername(username)
		am.AssertErrNil(t, err)
	})

	t.Run("invalid username", func(t *testing.T) {
		username := "no"
		err := validateUsername(username)
		am.AssertErrNotNil(t, err)
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		password := "correct horse battery staple"
		err := validatePassword(password)
		am.AssertErrNil(t, err)
	})

	t.Run("invalid password", func(t *testing.T) {
		password := "insecure"
		err := validatePassword(password)
		am.AssertErrNotNil(t, err)
	})
}

func TestValidateFullName(t *testing.T) {
	t.Run("valid full name", func(t *testing.T) {
		fullname := "abcdefghijklmnopqrstuvwx"
		err := validateFullName(fullname)
		am.AssertErrNil(t, err)
	})

	t.Run("invalid full name", func(t *testing.T) {
		fullname := "no"
		err := validateFullName(fullname)
		am.AssertErrNotNil(t, err)
	})
}

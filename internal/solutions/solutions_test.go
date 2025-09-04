package solutions

import (
	"testing"

	"github.com/mithcs/probox-api/pkg/am"
)

func TestValidateTitle(t *testing.T) {
	t.Run("valid title", func(t *testing.T) {
		title := "a valid title"

		err := validateTitle(title)
		am.AssertErrNil(t, err)
	})

	t.Run("invalid title", func(t *testing.T) {
		title := "no"

		err := validateTitle(title)
		am.AssertErrNotNil(t, err)
	})
}

func TestValidateDescription(t *testing.T) {
	t.Run("valid description", func(t *testing.T) {
		description := "a valid description"

		err := validateDescription(description)
		am.AssertErrNil(t, err)
	})

	t.Run("invalid description", func(t *testing.T) {
		description := "no"

		err := validateDescription(description)
		am.AssertErrNotNil(t, err)
	})
}

package solutions

import (
	"testing"

	"github.com/mithcs/probox-api/pkg/am"
)

func TestValidateSolution(t *testing.T) {
	t.Run("valid solution", func(t *testing.T) {
		solution := "a valid solution"

		err := validateSolution(solution)
		am.AssertErrNil(t, err)
	})

	t.Run("invalid solution", func(t *testing.T) {
		solution := "no"

		err := validateSolution(solution)
		am.AssertErrNotNil(t, err)
	})
}

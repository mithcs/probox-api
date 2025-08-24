package solutions

import (
	"testing"

	"github.com/mithcs/probox-api/pkg/am"
)

func TestValidateSolution(t *testing.T) {
	t.Run("valid solution", func(t *testing.T) {
		solution := CreateSolutionRequest{
			ProblemId: 1,
			Solution:  "this is a solution",
		}

		err := solution.validate()
		am.AssertErrNil(t, err)
	})

	t.Run("invalid problem id", func(t *testing.T) {
		solution := CreateSolutionRequest{
			ProblemId: 0,
			Solution:  "this is a solution",
		}

		err := solution.validate()
		am.AssertErrNotNil(t, err)
	})

	t.Run("invalid solution", func(t *testing.T) {
		solution := CreateSolutionRequest{
			ProblemId: 1,
			Solution:  "solution",
		}

		err := solution.validate()
		am.AssertErrNotNil(t, err)
	})
}

package problems

import (
	"testing"

	"github.com/mithcs/probox-api/pkg/am"
)

func TestValidateProblem(t *testing.T) {
	t.Run("valid title and description", func(t *testing.T) {
		problems := CreateProblemRequest{
			Title:       "problem",
			Description: "this is a problem",
		}

		err := problems.validate()
		am.AssertErrNil(t, err)
	})

	t.Run("invalid title", func(t *testing.T) {
		problems := CreateProblemRequest{
			Title:       "invalid title",
			Description: "this is a problem",
		}

		err := problems.validate()
		am.AssertErrNotNil(t, err)
	})

	t.Run("invalid description", func(t *testing.T) {
		problems := CreateProblemRequest{
			Title:       "title",
			Description: "invalid description",
		}

		err := problems.validate()
		am.AssertErrNotNil(t, err)
	})
}

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGitDir(t *testing.T) {
	t.Run("is git directory", func(t *testing.T) {
		gitClient := NewMockGitClient()

		gitClient.AddCommand(
			[]string{"rev-parse", "--is-inside-work-tree"},
			"true",
			nil,
		)

		ok, _ := IsGitDir(gitClient)

		assert.True(t, ok)
	})

	t.Run("is not git directory", func(t *testing.T) {
		gitClient := NewMockGitClient()

		gitClient.AddCommand(
			[]string{"rev-parse", "--is-inside-work-tree"},
			"fatal: not a git repository (or any of the parent directories): .git",
			nil,
		)

		ok, _ := IsGitDir(gitClient)

		assert.False(t, ok)
	})
}

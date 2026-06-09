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

		ok, err := IsGitDir(gitClient)

		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("output is false", func(t *testing.T) {
		gitClient := NewMockGitClient()

		gitClient.AddCommand(
			[]string{"rev-parse", "--is-inside-work-tree"},
			"false",
			nil,
		)

		ok, err := IsGitDir(gitClient)

		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("exit code 128: not a git repository", func(t *testing.T) {
		gitClient := NewMockGitClient()

		gitClient.AddExitError(
			[]string{"rev-parse", "--is-inside-work-tree"},
			"fatal: not a git repository (or any of the parent directories): .git",
			128,
		)

		ok, err := IsGitDir(gitClient)

		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("other exit codes: return error", func(t *testing.T) {
		gitClient := NewMockGitClient()

		gitClient.AddExitError(
			[]string{"rev-parse", "--is-inside-work-tree"},
			"some other error",
			1,
		)

		ok, err := IsGitDir(gitClient)

		assert.Error(t, err)
		assert.False(t, ok)
	})
}

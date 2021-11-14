package config

import (
	"testing"
)

func TestParseYaml(t *testing.T) {

	// TODO: improve this test
	t.Run("test YAML config parsing", func(t *testing.T) {
		s := `users:
  work: &work_user
    name: John Doe
    email: john.doe@corp.com
  personal:
    github_user: &personal_github_user
      name: JD42
      email: JD42@users.noreply.github.com
    gitlab_user: &personal_gitlab_user
      name: JD42
      email: 786972-JD42@users.noreply.gitlab.com

hosts:
  github.com:
    default: *personal_github_user
    overrides:
      - owner: corp
        user: *work_user

  gitlab.com:
    default: *personal_gitlab_user`

		_, err := parseYaml(s)
		if err != nil {
			t.Errorf("%+v", err)
		}
	})
}

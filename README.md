# zit

_git identity manager_

## How it works

_zit_ chooses a git identity based on:

1. git remote host
2. repository owner
3. repository name

… as defined in the configuration file:

```yaml
# Example config
---
users:
  work: &work_user
    name: "John Doe"
    email: "john.doe@corp.com"

  personal:
    github_user: &personal_github_user
      name: "JD42"
      email: "JD42@users.noreply.github.com"
      sign:
        key: "~/.ssh/id_ed25519_github.pub"
        format: "ssh"

    gitlab_user: &personal_gitlab_user
      name: "JD42"
      email: "786972-JD42@users.noreply.gitlab.com"

hosts:
  github.com:
    default: *personal_github_user
    overrides:
      - owner: "corp"
        user: *work_user

  gitlab.com:
    default: *personal_gitlab_user
```

## Setup

There are 2 ways to set up a configuration file:

1. Place it in
   [XDG_CONFIG_HOME](https://specifications.freedesktop.org/basedir-spec/0.6/)
   location: `$XDG_CONFIG_HOME/zit/config.yaml`
2. Place it in `.config` location: `$HOME/.config/zit/config.yaml`

Alternatively, run:

```sh
zit config init

vim $(zit config path)
```

## Usage

To set up an identity, run `zit set` inside a repo directory:

```bash
$ zit set  # personal repo
set user: jdoe <jdoe@users.noreply.github.com>
set signing key: ssh key at ~/.ssh/id_ed25519_github.pub

$ git remote get-url origin
https://github.com/jdoe/repo.git
```

```bash
$ zit set  # work repo
set user: John Doe <john.doe@corp.com>

$ git remote get-url origin
git@github.corp.com:team/repo.git
```

**Note**: Use `--dry-run` flag to test which identity will be used without
applying it.

## Installation

**On Mac/Linux with Homebrew**

```bash
brew tap ayakovlenko/tools
brew install ayakovlenko/tools/zit
```

**From sources**

```bash
git clone https://github.com/ayakovlenko/zit.git
cd zit
go install
```

**From binaries**

Download binaries from the
[releases](https://github.com/ayakovlenko/zit/releases) page.

## Setup

**Remove any existing global identity**

```bash
git config --unset-all --global user.name
git config --unset-all --global user.email
git config --unset-all --system user.name
git config --unset-all --system user.email
```

**Require config to exist in order to make commits**

```bash
git config --global user.useConfigOnly true
```

Without the global user name and user email, git would use the system's hostname
and username to make commits. Tell git to throw an error instead, requiring you
to specify an identity for every new project.

Run `zit doctor` to make sure the system is configured correctly:

```bash
$ zit doctor
- [x] git config --global user.useConfigOnly true
- [x] git config --unset-all --global user.name
- [x] git config --unset-all --global user.email
- [x] git config --unset-all --system user.name
- [x] git config --unset-all --system user.email
```

## Contributing

If you feel you can contribute to this project, or you've found a bug, create an issue or pull request.

This project is soley mantained, so it is prone to bugs and anti-patterns, please call them out.

All contributions are highly appreciated!

[![](https://api.star-history.com/svg?repos=ayakovlenko/zit&type=Date)](https://star-history.com/#ayakovlenko/zit&Date)

### Contributors

- [@arjunrn](https://github.com/arjunrn)
- [@etolmach](https://github.com/etolmach)

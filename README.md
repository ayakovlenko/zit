# zit

[![](https://api.codacy.com/project/badge/Grade/13955840a985457f8f2e5f22beeea75c)](https://www.codacy.com/manual/ayakovlenko/zit)

_git identity manager_

## How it works

_zit_ chooses a git identity based on (1) git remote host, (2) repository owner
name, (3) repository name as defined in a global configuration file
`$HOME/.zit/config.jsonnet`:

```jsonnet
local User(name, email) = { name: name, email: email };

local user = {
  "personal": User("jdoe", "jdoe@users.noreply.github.com"),
  "work": User("John Doe", "john.doe@corp.com")
};

{
  "github.com": { // Public GitHub
    "default": user.personal,
    "overrides": [
      { // Your employer's organization at github.com
        "owner": "corp",
        "user": user.work
      }
    ]
  },
  "github.corp.com": { // Your employer's GitHub Enterprise
    "default": user.work
  }
}
```

To set up an identity, run `zit` inside the repo:

```bash
$ zit  # personal repo
set user: jdoe <jdoe@users.noreply.github.com>

$ git remote get-url origin
https://github.com/jdoe/repo.git
```

```bash
$ zit  # work repo
set user: John Doe <john.doe@corp.com>

$ git remote get-url origin
git@github.corp.com:team/repo.git
```

## Installation

On Mac/Linux with Homebrew:

```bash
brew tap ayakovlenko/zit
brew install ayakovlenko/zit/zit
```

From sources:

```bash
git clone https://github.com/ayakovlenko/zit.git
cd zit
go install
```

From binaries:

Download binaries from the [releases](https://github.com/ayakovlenko/zit/releases) page.

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

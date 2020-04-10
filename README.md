# zit

_git identity manager_

## Installation

On Mac/Linux with Homebrew:

```
brew tap ayakovlenko/zit
brew install ayakovlenko/zit/zit
```

From sources:

```
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

Create a global configuration file at `$HOME/.zit/config.jsonnet`:

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
  "github.corp.com": { // GitHub Enterprise
    "default": user.work
  }
}
```

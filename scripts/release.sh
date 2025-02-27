#!/usr/bin/env sh
set -ea

# make sure the branch is main
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
  echo error: must be on 'main' branch
  exit 1
fi

VERSION=$(go run cmd/zit/main.go version)

git tag "$VERSION"

git push origin "$VERSION"

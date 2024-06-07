#!/usr/bin/env sh
set -eax

VERSION=$(go run cmd/zit/main.go version)
git tag "$VERSION"
git push origin "$VERSION"

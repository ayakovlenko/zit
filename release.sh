#!/usr/bin/env sh
set -eax

VERSION=$(go run . version)
git tag $VERSION
git push origin $VERSION

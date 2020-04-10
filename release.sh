#!/usr/bin/env sh
VERSION=$(go run . version)
git tag $VERSION
git push origin $VERSION

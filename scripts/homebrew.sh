#!/usr/bin/env sh
set -e

VERSION_TAG=$(go run cmd/zit/main.go version)     # vSEMVER
VERSION=$(echo "$VERSION_TAG" | sed -e "s/^v//g") #  SEMVER
FILE="zit-$VERSION.tar.gz"

cd ~/Downloads || exit

curl --fail -L "https://github.com/ayakovlenko/zit/archive/refs/tags/$VERSION_TAG.tar.gz" -o "$FILE"

echo "$FILE"

SUM=$(shasum -a 256 "$FILE" | awk '{ print $1; }')

echo "$SUM"

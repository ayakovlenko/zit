#!/usr/bin/env sh
set -e

MODULE=${MODULE:-"./..."}

if [ "$1" == "all" ]; then
    golangci-lint run --config .golangci.yaml "$MODULE"
    exit 0
fi

if [ "$1" == "only" ]; then
    golangci-lint run --no-config --enable "$2" "$MODULE"
    exit 0
fi

echo "Usage:"
echo
echo "  - all"
echo "  - only <linter>"

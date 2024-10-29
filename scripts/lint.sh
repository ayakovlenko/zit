#!/usr/bin/env bash
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

if [ "$1" == "linted" ]; then
	# The following packages are already linted and should stay this way.
	#
	# Should be replaced by simply `./lint all` in the future.
	packages=(
		'./cmd/zit/...'
		'./internal/app/...'
		'./pkg/xdg/...'
	)

	for package in "${packages[@]}"; do
		echo "linting package $package"
		golangci-lint run --config .golangci.yaml "$package"
	done

	exit 0
fi

echo "Usage:"
echo
echo "  - all"
echo "  - only <linter>"
echo "  - linted"

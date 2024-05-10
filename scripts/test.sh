#!/usr/bin/env sh
set -eax

go test -covermode=count -coverprofile=coverage.out ./...

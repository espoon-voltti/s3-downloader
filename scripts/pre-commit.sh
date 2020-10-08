#!/bin/sh

# Git pre-commit hook for linting Go source files
# See README for installation instructions

[ -z "$(gofmt -l .)" ] && exit 0

echo >&2 "Incorrectly formatted files. Run: make fmt"
exit 1

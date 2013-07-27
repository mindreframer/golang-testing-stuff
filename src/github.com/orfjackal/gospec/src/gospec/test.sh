#!/bin/sh
set -e

# Must disable inlining, or tests will fail because of stack traces
# not containing the innermost method call.
go test -gcflags='-l'

#!/usr/bin/env bash

# scripts/cibuild.sh: Setup environment for CI to run tests. This is primarily
#                 designed to run on the continuous integration server.

set -e

cd "$(dirname "$0")/.."

echo "Running tests …"
date "+%H:%M:%S"

# run tests.
scripts/test.sh
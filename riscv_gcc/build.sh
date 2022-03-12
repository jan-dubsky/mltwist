#!/bin/sh
set -eu

DIR="$(dirname -- "$0")"
readonly DIR
cd "$DIR"

readonly tag="${1-latest}"

docker build -t google2/riscv-gnu-toolchain-rv64ima:"$tag" .

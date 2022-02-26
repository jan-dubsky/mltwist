#!/bin/sh
set -eu

DIR="$(dirname -- "$0")"
cd "$DIR"

docker build -t riscv-gcc .

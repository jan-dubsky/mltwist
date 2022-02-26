#!/bin/sh
set -eu

usage(){
	cat <<EOF
Usage: compile.sh COMMANDS [[OPTIONS] ARGS]
EOF
}

if [ $# -lt 1 ]; then
	printf 'error: At least 1 source file must be provided.\n' >&2
	usage >&2
	exit 1
fi

DIR="$(pwd)"
readonly DIR

if printf '%s\n' "$DIR" | grep -q :; then
	printf "error: Build directory cannot contain colon in it: '%s'.\n" "$DIR" >&2
	exit 1
fi

docker run --rm -v "$DIR:/build" riscv-gcc "$@"

#!/bin/sh
set -eu

readonly DOCKER_BASE=google2/riscv-gnu-toolchain-rv64ima

usage() {
	cat <<EOF
Usage: build.sh TAG
       build.sh -h

Options:
  -h  Print this help.
EOF
}

DIR="$(dirname -- "$0")"
readonly DIR
cd "$DIR"

while getopts h opt; do
	case $opt in
	h) usage; exit ;;
	*) usage >&2; exit 1 ;;
	esac
done
shift $((OPTIND - 1))

if [ $# -ne 1 ]; then
	usage >&2
	exit 1
fi

readonly TAG="$1"
shift

if
	! expr "$TAG" : '[0-9A-Za-z_][0-9A-Za-z_.-]*$' >/dev/null || \
	[ ${#TAG} -gt 128 ]
then
	printf 'Invalid tag provided: %s\n' "$TAG" >&2
	exit 1
fi

readonly DOCKER_IMAGE="$DOCKER_BASE:$TAG"

docker build -t "$DOCKER_IMAGE" .
docker image tag "$DOCKER_IMAGE" "$DOCKER_BASE:latest"

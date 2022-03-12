#!/bin/sh
set -eu

local_dir="$(dirname -- "$0")"
dir="$(cd "$local_dir" && pwd)"
readonly dir

readonly repodir="$dir/git_src"

if ! [ -d "$repodir" ]; then
	git clone --depth=1 https://git.savannah.gnu.org/git/grep.git "$repodir"
fi

chmod -R o+w "$repodir"
cd "$repodir"

git checkout master
git pull --ff-only

cmd='
	set -eu
	set -x		# Is useful for debugging of build problems.

	./bootstrap
	./configure --host=riscv64-unknown-linux-gnu
	make CFLAGS=-Wno-cast-align all
'
"$dir"/../../riscv_gcc/run.sh /bin/sh -c "$cmd"

# We cannot use cp as we need to change owner from the in-container user to
# current user.
cat "$repodir"/src/grep > "$dir"/grep
chmod +x "$dir"/grep

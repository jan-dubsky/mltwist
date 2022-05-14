
GIT_DIR=git_src
GREP_BUILD=$(GIT_DIR)/src/grep
BUILD=grep

build: grep FORCE

$(BUILD): $(GREP_BUILD)
	# We cannot use cp as we need to change owner from the in-container
	# user to current user.
	cat $< > $@
	chmod +x $@

$(GREP_BUILD): $(GIT_DIR)
	# For some reason the riscv GCC doesn't like unaligned memory accesses
	# even though RISC-V supports them. The
	# -Wno-cast-align is necessary for the build not to fail.
	
	cd $(GIT_DIR) && ../../../riscv_gcc/run.sh /bin/sh -c ' \
		set -eu;							\
		set -x;								\
										\
		./bootstrap;							\
		./configure --host=riscv64-unknown-linux-gnu LDFLAGS=-static;   \
		make CFLAGS=-Wno-cast-align all; 				\
	'

$(GIT_DIR):
	git clone --depth=1 https://git.savannah.gnu.org/git/grep.git $@
	chmod -R o+w $@

clean:
	rm -rf $(GIT_DIR)
	rm -f $(BUILD)

FORCE:

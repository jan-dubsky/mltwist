
GIT_DIR=git_src
GREP_BUILD=$(GIT_DIR)/src/grep
BUILD=grep

build: $(BUILD) FORCE

$(BUILD): $(GREP_BUILD)
	# We cannot use cp as we need to change owner from the in-container
	# user to current user.
	cat $< > $@
	chmod +x $@

$(GIT_DIR)/configure: $(GIT_DIR)/bootstrap
	cd $(GIT_DIR) && ../../../riscv_gcc/run.sh /bin/sh -c ' \
		set -eu;                                        \
		./bootstrap;                                    \
	'

$(GIT_DIR)/Makefile: $(GIT_DIR)/configure
	cd $(GIT_DIR) && ../../../riscv_gcc/run.sh /bin/sh -c ' \
		set -eu;                                        \
		./configure --host=riscv64-unknown-linux-gnu    \
	'


$(GREP_BUILD): $(GIT_DIR)/Makefile
	# For some reason the riscv GCC doesn't like unaligned memory accesses
	# even though RISC-V supports them. The
	# -Wno-cast-align is necessary for the build not to fail.
	#
	# The --whole-archive bypass with -lpthread is necessary for the
	# compiler not tp produce jumpts to address zero:
	# https://gcc.gnu.org/bugzilla/show_bug.cgi?id=33960#c9
	
	cd $(GIT_DIR) && ../../../riscv_gcc/run.sh /bin/sh -c ' \
		set -eu;                                        \
		                                                \
		make CFLAGS=-Wno-cast-align all;                \
	'

$(GIT_DIR)/bootstrap:
	git clone --depth=1 https://git.savannah.gnu.org/git/grep.git $@
	chmod -R o+w $@

prune:
	rm -rf $(GIT_DIR)

clean: prune
	rm -f $(BUILD)

FORCE:

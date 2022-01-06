
BUILD_DIR=build

CMDS=$(shell ls cmd)
BINS=$(addprefix $(BUILD_DIR)/,$(CMDS))

build: $(BINS)

$(BUILD_DIR)/%: cmd/% FORCE
	@printf "Compiling %s...\n" $@
	@mkdir -p $(BUILD_DIR)
	@go build -o $@ ./$<

test: build FORCE
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

FORCE:

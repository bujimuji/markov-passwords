# The binary to build
BIN ?= mkpass

# This repo's root import path
PKG := github.com/bujimuji/markov-passwords

all:
	@$(MAKE) mkpass

mkpass: bin/$(BIN)

bin/$(BIN): build-dirs
	@echo "building: $@"
	go build \
	-o bin/$(BIN) \
	$(PKG)/cmd/markov-passwords

build-dirs:
	@mkdir -p bin

clean:
	rm -rf bin
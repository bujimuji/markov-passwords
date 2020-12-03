# The binary to build
TRAIN ?= train
SAMPLE ?= sample

# This repo's root import path
PKG := github.com/bujimuji/markov-passwords

all:
	@$(MAKE) train
	@$(MAKE) sample

train: bin/$(TRAIN)
sample: bin/$(SAMPLE)

bin/$(TRAIN): build-dirs
	@echo "building: $@"
	go build \
	-o bin/$(TRAIN) \
	$(PKG)/cmd/$(TRAIN)

bin/$(SAMPLE): build-dirs
	@echo "building: $@"
	go build \
	-o bin/$(SAMPLE) \
	$(PKG)/cmd/$(SAMPLE)

build-dirs:
	@mkdir -p bin

clean:
	rm -rf bin
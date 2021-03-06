MODULE?=agent
BUILD_TAGS?=$(MODULE)
BUILD_FLAGS=-ldflags "-X dudu/version.GitCommit=`git rev-parse --short=8 HEAD`"

########################### build ######################
build:
	go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o build/dudu-$(MODULE) ./cmd/$(MODULE)/

build_race:
	go build -race $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o build/dudu-$(MODULE) ./cmd/$(MODULE)/

install:
	go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o $(GOPATH)/bin/dudu-$(MODULE) ./cmd/$(MODULE)/

clean:
	rm -rf build/

.PHONY: build build_race install

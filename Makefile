PROJECT := github.com/redhat-developer/kdo
GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
PKGS := $(shell go list  ./... | grep -v $(PROJECT)/vendor | grep -v $(PROJECT)/tests )
COMMON_FLAGS := -X $(PROJECT)/pkg/kdo/cli/version.GITCOMMIT=$(GITCOMMIT)
BUILD_FLAGS := -ldflags="-w $(COMMON_FLAGS)"
DEBUG_BUILD_FLAGS := -ldflags="$(COMMON_FLAGS)"
FILES := odo dist
TIMEOUT ?= 1800s

default: bin

.PHONY: bin
bin:
	go build ${BUILD_FLAGS} cmd/kdo/kdo.go

.PHONY: install
install:
	go install ${BUILD_FLAGS} ./cmd/kdo/

.PHONY: clean
clean:
	@rm -rf $(FILES)

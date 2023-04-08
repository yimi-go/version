# Build all by default, even if it's not first
.DEFAULT_GOAL := help

# ==============================================================================
# Build options
ROOT_PACKAGE=github.com/yimi-go/version
# Minimum test coverage
ifeq ($(origin COVERAGE),undefined)
COVERAGE := 60
endif

# ==============================================================================
# Common rules
SHELL := /bin/bash
SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(SELF_DIR) && pwd -P))
endif
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
endif
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
endif
# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*' 2> /dev/null)
endif
# Check if the tree is dirty.  default to dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2> /dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD 2> /dev/null)
FIND := find . ! -name '*.pb.go'
SEDCMD=$(shell which sed)
ifeq ($(shell uname),Darwin)
  SEDCMDI=$(SEDCMD) -i ''
  XARGS := xargs -r
  AWK := gawk
else
  SEDCMDI=$(SEDCMD) -i
  XARGS := xargs --no-run-if-empty
  AWK := awk
endif
# Makefile settings
ifndef V
MAKEFLAGS += --no-print-directory
endif
COMMA := ,
SPACE :=
SPACE +=
GO := go

# ==============================================================================
# Tools rules
.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.golangci-lint
install.golangci-lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) install github.com/jstemmer/go-junit-report@latest

.PHONY: install.gsemver
install.gsemver:
	@$(GO) install github.com/arnaud-deprez/gsemver@latest

.PHONY: install.git-chglog
install.git-chglog:
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

.PHONY: install.golines
install.golines:
	@$(GO) install github.com/segmentio/golines@latest

.PHONY: install.go-mod-outdated
install.go-mod-outdated:
	@$(GO) install github.com/psampaz/go-mod-outdated@latest

.PHONY: install.mockgen
install.mockgen:
	@$(GO) install github.com/golang/mock/mockgen@latest

.PHONY: install.gotests
install.gotests:
	@$(GO) install github.com/cweill/gotests/gotests@latest

.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	@$(GO) install github.com/golang/protobuf/protoc-gen-go@latest

.PHONY: install.protoc-gen-validate
install.protoc-gen-validate:
	@mkdir -p "$(TMP_DIR)/pgv.install"
	@echo "cloning into $(TMP_DIR)/pgv.install" \
		&& git clone https://github.com/bufbuild/protoc-gen-validate.git "$(TMP_DIR)/pgv.install"
	@echo "install protoc-gen-validate from source" \
		&& cd "$(TMP_DIR)/pgv.install" && make build

.PHONY: install.protoc-gen-go-grpc
install.protoc-gen-go-grpc:
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: install.protoc-gen-connect-go
install.protoc-gen-connect-go:
	@$(GO) install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

.PHONY: install.grpcurl
install.grpcurl:
	@$(GO) install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

.PHONY: install.goimports
install.goimports:
	@(GO) install golang.org/x/tools/cmd/goimports@latest

.PHONY: install.go-gitlint
install.go-gitlint:
	@$(GO) install github.com/llorllale/go-gitlint/cmd/go-gitlint@latest

# ==============================================================================
# Golang rules
ifneq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
endif
ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif
GOPATH := $(shell $(GO) env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif
EXCLUDE_TESTS=$(shell cat $(ROOT_DIR)/.gotestignore)

.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

.PHONY: go.format
go.format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formatting codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .
	@go mod edit -fmt

ifeq ($(origin V), 1)
GO_TEST_FLAG += -v
endif

.PHONY: go.test
go.test: tools.verify.go-junit-report
	@mkdir -p $(OUTPUT_DIR)
	@echo "===========> Run unit test"
	@$(GO) test $(GO_BUILD_FLAGS) $(GO_TEST_FLAG) $(ROOT_PACKAGE)/...
	@set -o pipefail;$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short $(GO_TEST_FLAG) `go list ./...|\
		egrep -v $(subst $(SPACE),'|',$(sort $(EXCLUDE_TESTS)))` 2>&1 | \
		tee >(go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml)
	@$(SEDCMDI) '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out # remove mock_.*.go files from test coverage
	@$(SEDCMDI) '/.*.pb.go/d' $(OUTPUT_DIR)/coverage.out # remove .*.pb.go files from test coverage
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

.PHONY: go.test.cover
go.test.cover: go.test
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
		$(AWK) -v target=$(COVERAGE) -f $(ROOT_DIR)/build/coverage.awk

.PHONY: go.updates
go.updates: tools.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  COVERAGE         Minimum test coverage. Default is 60.
  V                Set to 1 enable verbose build. Default is undefined.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## git.hooks: Setup the dev environment.
.PHONY: git.hooks
git.hooks: tools.verify.go-gitlint
	@cmp -s build/githooks/commit-msg .git/hooks/commit-msg \
		&& cmp -s build/githooks/pre-commit .git/hooks/pre-commit \
		|| echo 'clone git hooks...' \
		&& cp -f build/githooks/* .git/hooks/
	@chmod +x .git/hooks/*

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

## cover: Run unit test and get test coverage.
.PHONY: cover
cover:
	@$(MAKE) go.test.cover

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format:
	@$(MAKE) go.format

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

.PHONY: tidy
tidy:
	@$(GO) mod tidy

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"

.PHONY: release.tag
release.tag: tools.verify.gsemver release.ensure-tag
	@git push origin `git describe --tags --abbrev=0`

.PHONY: release.ensure-tag
release.ensure-tag: tools.verify.gsemver
	@build/scripts/ensure_tag.sh
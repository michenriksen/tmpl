# Makefile adapted from https://github.com/thockin/go-build-template
# Released under the Apache-2.0 license.
DBG_MAKEFILE ?=
ifeq ($(DBG_MAKEFILE),1)
    $(warning ***** starting Makefile for goal(s) "$(MAKECMDGOALS)")
    $(warning ***** $(shell date))
else
    # If we're not debugging the Makefile, don't echo recipes.
    MAKEFLAGS += -s
endif

DBG ?=

# We don't need make's built-in rules.
MAKEFLAGS += --no-builtin-rules
# Be pedantic about undefined variables.
MAKEFLAGS += --warn-undefined-variables

OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

GOBIN ?= "go"
GOFLAGS ?=
HTTP_PROXY ?=
HTTPS_PROXY ?=

default: verify

test: # @HELP run tests
test:
	./scripts/test.sh

autotest: # @HELP run tests on file changes
autotest:
	./scripts/autotest.sh

cover: # @HELP run tests with coverage and open HTML report
cover:
	./scripts/cover.sh

lint: # @HELP run all linting
lint: lint-go lint-yaml lint-schema lint-ci lint-docs

lint-go: # @HELP run Go linting
lint-go:
	./scripts/lint-go.sh

lint-yaml: # @HELP run YAML linting
lint-yaml:
	./scripts/lint-yaml.sh

lint-schema: # @HELP run JSON schema linting
lint-schema:
	./scripts/lint-schema.sh

lint-ci: # @HELP run linting on CI workflow files
lint-ci:
	./scripts/lint-ci.sh

lint-docs: # @HELP run documentation linting
lint-docs:
	./scripts/lint-docs.sh

verify: # @HELP run all tests and linters (default target)
verify: test lint

build: # @HELP build a snapshot binary for current OS and ARCH in ./dist/
build:
	./scripts/build.sh

next-version: # @HELP determine the next version for a release
next-version:
	${GOBIN} run github.com/caarlos0/svu@latest next

gen-docs: # @HELP generate documentation
gen-docs:
	./scripts/gen-docs.sh

reset-golden: # @HELP remove all generated golden test files
reset-golden:
	./scripts/reset-golden.sh

help: # @HELP print this message
help:
	echo "VARIABLES:"
	echo "  OS       = $(OS)"
	echo "  ARCH     = $(ARCH)"
	echo "  GOBIN    = $(GOBIN)"
	echo "  GOFLAGS  = $(GOFLAGS)"
	echo "  DBG      = $(DBG)"
	echo
	echo "TARGETS:"
	grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST)     \
	    | awk '                                   \
	        BEGIN {FS = ": *# *@HELP"};           \
	        { printf "  %-20s$ %s\n", $$1, $$2 };  \
	    '

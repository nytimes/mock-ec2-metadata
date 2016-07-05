#
#  Makefile
#
#  The kickoff point for all project management commands.
#

# Program version
VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/')

# Binary name for bintray
BIN_NAME=mock-ec2-metadata

# Project owner for bintray
OWNER=NYTimes

# Project name for bintray
PROJECT_NAME=mock-ec2-metadata

# Project url used for builds
# examples: github.com, bitbucket.org
REPO_HOST_URL=github.com

# Grab the current commit
GIT_COMMIT=$(shell git rev-parse HEAD)

# Check if there are uncommited changes
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# Use a local vendor directory for any dependencies; comment this out to
# use the global GOPATH instead
GOPATH=$(PWD)/.vendor

INSTALL_PATH=$(GOPATH)/src/github.com/NYTimes/mock-ec2-metadata

default: build

help:
	@echo 'Management commands for mock-ec2-metadata:'
	@echo
	@echo 'Usage:'
	@echo '    make build    Compile the project.'
	@echo '    make link     Symlink this project into the GOPATH.'
	@echo '    make test     Run tests on a compiled project.'
	@echo '    make fmt      Reformat the source tree with gofmt.'
	@echo '    make clean    Clean the directory tree.'
	@echo

build: .git $(GOPATH)/bin/gogpm $(INSTALL_PATH)
	@echo "building ${OWNER} ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	$(GOPATH)/bin/gogpm install &&\
	go build -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY}" -o bin/${BIN_NAME} main.go

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

.git:
	git init
	git add -A .
	git commit -m 'Initial scaffolding.'

link:
	# relink into the go path
	if [ ! $(INSTALL_PATH) -ef . ]; then \
		mkdir -p `dirname $(INSTALL_PATH)`; \
		ln -s $(PWD) $(INSTALL_PATH); \
	fi

$(INSTALL_PATH):
	make link

$(GOPATH)/bin/gogpm:
	go get github.com/mtibben/gogpm

test:
	go test ./...

fmt:
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

.PHONY: build dist clean test help default link fmt

#
#  Makefile
#
#  The kickoff point for all project management commands.
#

# Program version
VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/')
GOPATH := $(shell go env GOPATH)

default: build

help:
	@echo 'Management commands for mock-ec2-metadata:'
	@echo
	@echo 'Usage:'
	@echo '    make build    Compile the project.'
	@echo '    make deps     Download dependencies.'
	@echo '    make test     Run tests on a compiled project.'
	@echo '    make fmt      Reformat the source tree with gofmt.'
	@echo '    make clean    Clean the directory tree.'
	@echo

deps:
	@go get github.com/mitchellh/gox
	@go get github.com/tcnksm/ghr
	@go get -v .;
	@go mod tidy;

# Runs the go vet command, will be a dependency for any test.
vet:
	@go vet .
	@go vet ./...

build:
	$(GOPATH)/bin/gox -verbose -output "bin/mock-ec2-metadata_${VERSION}_{{.OS}}_{{.Arch}}" -os="linux" -os="darwin" -arch="386" -arch="amd64" -arch="arm64" ./cmd/server

release:
	$(GOPATH)/bin/ghr --username NYTimes --token ${GITHUB_TOKEN} -r mock-ec2-metadata --replace ${VERSION} bin/

clean:
	@test ! -e bin/|| rm -f bin/*

tests:
	go test ./...

fmt:
	find . -name '*.go' -exec gofmt -w=true {} ';'

.PHONY: build dist clean test help default fmt

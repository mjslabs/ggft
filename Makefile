PKG := $(shell head -1 go.mod | sed -e 's/module //')
PKG_LIST := $(shell go list ${PKG}/...)
GO_FILES := $(shell find . -name '*.go')
OUT := $(shell basename ${PKG} | sed -e 's/\.git//')
DIST := $(OUT).zip
VERSION := $(shell git describe --tag --long --dirty 2>/dev/null)
ifeq ($(VERSION),)
VERSION = development
endif

INSTALLBIN := $(GOPATH)/bin
ifeq ($(GOPATH),)
INSTALLBIN = ~/go/bin
endif

all: build

# Run all build steps
.PHONY: build
build: vet lint semgrep ggft

# Create binary
ifeq ($(TRAVIS_OS_NAME),windows)
OUT := $(OUT).exe
.PHONY: ggft
endif
ggft:
	go build -v -o $(OUT) -ldflags="-X 'github.com/mjslabs/ggft/cmd.version=$(VERSION)'" $(PKG)

# Package up in a release
ifeq ($(DIST_NAME),)
DIST_NAME = $(OUT)
endif
.PHONY: dist
dist: build
	zip $(DIST_NAME).zip $(OUT)

# Install locally
.PHONY: install
install: ggft
	install -D $(OUT) $(INSTALLBIN)

# Create test file
c.out:
	go test -coverprofile=c.out -v $(PKG)/...
	@echo Total coverage: `go tool cover -func c.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'`%

# Create coverage report in HTML
cover.html: c.out
	go tool cover -html=c.out -o cover.html

# Delete all artifacts
.PHONY: clean
clean:
	rm -f $(OUT) $(OUT).exe $(OUT).exe.zip $(DIST) c.out

# Run go linter
.PHONY: lint
lint:
	@for file in $(GO_FILES); do \
		golint $$file ; \
	done

# Run semgrep static code analysis
.PHONY: semgrep
semgrep:
	@docker run --rm -v "${PWD}:/src" returntocorp/semgrep --config=https://semgrep.dev/p/r2c-CI --exclude "*_test.go" .

# Run go test
.PHONY: test
test: c.out

# Tidy go.mod and .sum file
.PHONY: tidy
tidy:
	@go mod tidy

# Run go vet
.PHONY: vet
vet:
	@go vet $(PKG_LIST)

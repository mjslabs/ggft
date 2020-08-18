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

.PHONY: build
build: vet lint ggft


ifeq ($(GOOS),windows)
OUT := $(OUT).exe
.PHONY: ggft
endif
ggft:
	go build -v -o $(OUT) -ldflags="-X 'github.com/mjslabs/ggft/cmd.version=$(VERSION)'" $(PKG)

.PHONY: dist
dist: ggft
	zip $(OUT).zip $(OUT)

.PHONY: install
install: ggft
	install -D $(OUT) $(INSTALLBIN)

c.out:
	go test -coverprofile=c.out -v $(PKG)/...

.PHONY: clean
clean:
	rm -f $(OUT) $(OUT).exe $(OUT).exe.zip $(DIST) c.out

.PHONY: lint
lint:
	@for file in $(GO_FILES); do \
		golint $$file ; \
	done

.PHONY: test
test: c.out

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: vet
vet:
	@go vet $(PKG_LIST)

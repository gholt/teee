NAME=teee
TOP_IMPORT=github.com/gholt/$(NAME)
GIT_TAG=$(shell git describe --tags)
GIT_DIRTY=$(shell git status --porcelain)
VERSION=$(GIT_TAG)
ifneq "$(GIT_DIRTY)" ""
    VERSION=$(GIT_TAG)-dirty
endif

default: fmt vet test install

fmt:
	go fmt ./...

install:
	go install -ldflags "-X $(TOP_IMPORT).Version=$(VERSION)" ./...

install-race:
	go install -race -ldflags "-X $(TOP_IMPORT).Version=$(VERSION)" ./...

test:
	go test ./...

test-race:
	go test -race -cpu=1,2,7 ./...

vet:
	go vet ./...

build: fmt vet test build-all

build-all: build-linux build-windows build-macos

build-linux: FORCE
	mkdir -p build
	GOOS=linux go build -o build/$(NAME)-$(VERSION)-linux -ldflags "-X $(TOP_IMPORT).Version=$(VERSION)" $(TOP_IMPORT)/$(NAME)

build-windows: FORCE
	mkdir -p build
	GOOS=windows go build -o build/$(NAME)-$(VERSION)-windows.exe -ldflags "-X $(TOP_IMPORT).Version=$(VERSION)" $(TOP_IMPORT)/$(NAME)

build-macos: FORCE
	mkdir -p build
	GOOS=darwin go build -o build/$(NAME)-$(VERSION)-macos -ldflags "-X $(TOP_IMPORT).Version=$(VERSION)" $(TOP_IMPORT)/$(NAME)

FORCE:

clean:
	rm -rf build

build_commit=$(shell git rev-parse HEAD)
build_version=$(shell git describe --tags 2> /dev/null || echo "dev-$(shell git rev-parse HEAD)")

GOBUILDPATH ?= $(CURDIR)
GOTESTFLAGS := 

DESTDIR ?= $(CURDIR)/dist

.PHONY: all
all: deps
all: build
all: install

.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: build
build:
	go build -o $(GOBUILDPATH)/build/lmu $(CURDIR)/cmd/main.go

.PHONY: install
install: | $(DESTDIR)
	cp $(CURDIR)/configs/defaultConf.yaml $(DESTDIR)/etc/lmuConf.yaml
	cp $(GOBUILDPATH)/build/lmu $(DESTDIR)/usr/bin/lmu

$(DESTDIR):
	mkdir $(DESTDIR)
	mkdir $(DESTDIR)/etc
	mkdir $(DESTDIR)/usr
	mkdir $(DESTDIR)/usr/bin

.PHONY: test
test:
	go test $(GOTESTFLAGS) $(CURDIR)/internal/tests

.PHONY: clean
clean:
	rm -r build || true
	rm -r dist|| true

.PHONY: version
version:
	 @echo "Version:           ${build_version}"
	 @echo "Git Commit:        ${build_commit}"

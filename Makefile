build_commit=$(shell git rev-parse HEAD)
build_version=$(shell git describe --tags 2> /dev/null || echo "dev-$(shell git rev-parse HEAD)")

GOBUILDPATH ?= $(CURDIR)
GOTESTFLAGS := 

DESTDIR ?= 

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
install: install_config install_binary
	
	

.PHONY: install_config
install_config:| $(DESTDIR)/etc/lmu 
	cp $(CURDIR)/configs/defaultConfig.yaml $(DESTDIR)/etc/lmu/lmuConfig.yaml
	cp $(CURDIR)/configs/defaultSyscalls.yaml $(DESTDIR)/etc/lmu/lmuSyscalls.yaml

.PHONY: install_binary
install_config:| $(DESTDIR)/usr/bin
	cp $(GOBUILDPATH)/build/lmu $(DESTDIR)/usr/bin/lmu


$(DESTDIR)/etc/lmu:
	mkdir -p $(DESTDIR)/etc/lmu
	

$(DESTDIR)/usr/bin:
	mkdir -p $(DESTDIR)/usr/bin

.PHONY: test
test:
	go test $(GOTESTFLAGS) $(CURDIR)/internal/tests

.PHONY: clean
clean:
	rm -r build || true
	
.PHONY: version
version:
	 @echo "Version:           ${build_version}"
	 @echo "Git Commit:        ${build_commit}"

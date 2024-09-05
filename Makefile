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
	/usr/local/go/bin/go mod download
	/usr/local/go/bin/go mod tidy

.PHONY: build
build:
	/usr/local/go/bin/go build -o $(GOBUILDPATH)/build/lmu $(CURDIR)/cmd/main.go

.PHONY: install
install: install_config install_binary
	
.PHONY: install_config
install_config:| $(DESTDIR)/etc/lmu 
	cp $(CURDIR)/configs/defaultConfig.yaml $(DESTDIR)/etc/lmu/lmuConfig.yaml

.PHONY: install_binary
install_binary:| $(DESTDIR)/usr/bin
	cp $(GOBUILDPATH)/build/lmu $(DESTDIR)/usr/bin/lmu


$(DESTDIR)/etc/lmu:
	mkdir -p $(DESTDIR)/etc/lmu
	

$(DESTDIR)/usr/bin:
	mkdir -p $(DESTDIR)/usr/bin

.PHONY: test
test:
	/usr/local/go/bin/go test $(GOTESTFLAGS) $(CURDIR)/internal/tests

.PHONY: clean
clean:
	rm -r build || true
	
.PHONY: version
version:
	 @echo "Version:           ${build_version}"
	 @echo "Git Commit:        ${build_commit}"

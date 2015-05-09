PREFIX ?= /usr/local
SRC = $(wildcard *.go)
BUILD_NAME=nodester
export GOPATH=$(CURDIR)/Godeps/_workspace


nodester:
	@mkdir build
	go build -o build/$(BUILD_NAME) $(SRC)

install: nodester
	install -m 0755 build/$(BUILD_NAME) $(PREFIX)/bin/$(BUILD_NAME)

uninstall:
	@rm $(PREFIX)/bin/$(BUILD_NAME)

clean:
	@rm -rf build

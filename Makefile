VERSION := $(or $(VERSION),$(shell git describe --tags --long --dirty --always))
LDFLAGS := -s -w -X github.com/tomzxcode/gh-cached/internal/version.Version=$(VERSION)
OUTPUT ?= gh-cached

.PHONY: build clean

build:
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(OUTPUT) .

clean:
	rm -f $(OUTPUT)

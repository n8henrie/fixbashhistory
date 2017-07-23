SHELL := /bin/bash
PWD = $(shell pwd)
GREP := $(shell command -v ggrep || command -v grep)
LDFLAGS := -ldflags "-X main.version=$(git describe --tags)"

help:
	@$(GREP) --only-matching --word-regexp '^[^[:space:].]*:' Makefile | sed 's|:[:space:]*||'

test:
	go test -v .

bench:
	go test -bench=.

rpi:
	GOOS=linux GOARCH=arm GOARM=7 go build $(LDFLAGS) -v .

macos:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -v .

clean:
	-rm fixbashhistory

.PHONY: clean help

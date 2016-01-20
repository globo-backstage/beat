include go.mk

.PHONY: build
build: gomkbuild

.PHONY: xbuild
xbuild: gomkxbuild

.PHONY: clean
clean: gomkclean

.PHONY: run
run:
	go run beat/main.go

.PHONY: setup
setup: deps restoregodeps

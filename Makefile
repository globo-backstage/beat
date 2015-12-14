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

save-deps:
	$(GOPATH)/bin/godep save ./...

setup:
	go get $(GO_EXTRAFLAGS) -u -d -t ./...
	go get $(GO_EXTRAFLAGS) github.com/tools/godep
	$(GOPATH)/bin/godep restore ./...
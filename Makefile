setup:
	go get $(GO_EXTRAFLAGS) -u -d -t ./...
	go get $(GO_EXTRAFLAGS) github.com/tools/godep
	$(GOPATH)/bin/godep restore ./...

run:
	go run beat/main.go

save-deps:
	$(GOPATH)/bin/godep save ./...

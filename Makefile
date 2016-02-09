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

.PHONY: doc-server
doc-server:
	mkdocs serve

.PHONY: setup-docs
setup-docs:
	pip install -r requirements_docs.txt

.PHONY: deploy-docs
deploy-docs:
	mkdocs gh-deploy --clean

.PHONY: update-mocks
update-mocks:
	mockgen -destination "mocks/mock_db/mock_db.go" github.com/backstage/beat/db Database

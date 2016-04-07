.PHONY: all setup run

all: setup run

setup:
	npm install

run:
	node index.js

post:
	curl -H "Content-Type: application/json" -XPOST http://localhost:8100 -d'{"hook": "beforeSave", "customCodes": ["scheduler"]}'

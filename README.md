[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/backstage/beat?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/backstage/beat.png?branch=master)](https://travis-ci.org/backstage/beat)

## What is Backstage Beat?

Backstage Beat is a Backend-as-a-Service software, that makes mobile and web development to be really fast with restful APIs.

**Note: The software is still in hard development, not ready for production environments**


## Getting started

### Requirements

- Go 1.5+
- MongoDB 3+

## Download and install the devolpement version

Ensure if your GOPATH environment variable is setted, see more in: https://golang.org/doc/code.html#GOPATH

```
go get "github.com/backstage/beat/beat"
cd $GOPATH/src/github.com/backstage/beat
make setup
```

## Run the devolpement version

```
make run
```


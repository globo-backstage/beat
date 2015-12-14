package main

import (
	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/server"
)

func main()  {
	s := server.New(&auth.DraftAuthentication{})
	s.Run()
}

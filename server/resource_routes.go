package server

import (
	"net/http"
	"fmt"
)

func (s *Server) createResource(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: Create resource")
}

func (s *Server) findResource(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: Find resource")
}

func (s *Server) findOneResource(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: findOne resource")
}

func (s *Server) findResourceById(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: Find resource by id")
}

func (s *Server) deleteResourceById(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: delete resource By Id")
}

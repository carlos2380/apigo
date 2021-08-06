package server

import (
	"apigo/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router http.Handler
}

func NewServer() *Server {
	r := mux.NewRouter()
	r.HandleFunc("/api/stores", handlers.GetStore).Methods("GET")

	return &Server{
		Router: r,
	}
}

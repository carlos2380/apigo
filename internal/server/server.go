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

	r.HandleFunc("/api/stores", handlers.GetStores).Methods(http.MethodGet)
	r.HandleFunc("/api/stores/{id}", handlers.GetStore).Methods(http.MethodGet)
	r.HandleFunc("/api/stores/{id}", handlers.DeleteStore).Methods(http.MethodDelete)
	r.HandleFunc("/api/stores", handlers.PostStore).Methods(http.MethodPost)
	r.HandleFunc("/api/stores/{id}", handlers.PutStore).Methods(http.MethodPut)

	return &Server{
		Router: r,
	}
}

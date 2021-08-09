package server

import (
	"apigo/internal/handlers"
	"apigo/internal/storage/postgres"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router http.Handler
}

func NewServer(ctrlDB *postgres.PostgresDB) *Server {
	r := mux.NewRouter()
	env := &handlers.EnvHandler{CtrlDB: ctrlDB}
	env.CtrlDB.InitDB()
	r.HandleFunc("/api/stores", env.GetStores).Methods(http.MethodGet)
	r.HandleFunc("/api/stores/{id}", env.GetStore).Methods(http.MethodGet)
	r.HandleFunc("/api/stores/{id}", env.DeleteStore).Methods(http.MethodDelete)
	r.HandleFunc("/api/stores", env.PostStore).Methods(http.MethodPost)
	r.HandleFunc("/api/stores/{id}", env.PutStore).Methods(http.MethodPut)

	return &Server{
		Router: r,
	}
}

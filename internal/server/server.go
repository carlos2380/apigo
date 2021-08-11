package server

import (
	"apigo/internal/handlers"
	"apigo/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(ctrlDB storage.DataBase) http.Handler {
	r := mux.NewRouter()
	env := &handlers.EnvHandler{CtrlDB: ctrlDB}
	// env.CtrlDB.InitDB()
	r.HandleFunc("/api/stores", env.GetStores).Methods(http.MethodGet)
	r.HandleFunc("/api/stores/{id}", env.GetStore).Methods(http.MethodGet)
	r.HandleFunc("/api/stores/{id}", env.DeleteStore).Methods(http.MethodDelete)
	r.HandleFunc("/api/stores", env.PostStore).Methods(http.MethodPost)
	r.HandleFunc("/api/stores/{id}", env.PutStore).Methods(http.MethodPut)

	return r
}

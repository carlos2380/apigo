package server

import (
	"apigo/internal/handlers"
	"apigo/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(stgStore storage.StoreStorage, stgCustomer storage.CustomerStorage) http.Handler {
	r := mux.NewRouter()
	stgHandler := &handlers.StorageHandler{StgStore: stgStore, StgCustomer: stgCustomer}
	r.HandleFunc("/api/stores", stgHandler.GetStores).Methods(http.MethodGet)
	r.HandleFunc("/api/stores/{id}", stgHandler.GetStore).Methods(http.MethodGet)
	r.HandleFunc("/api/stores/{id}", stgHandler.DeleteStore).Methods(http.MethodDelete)
	r.HandleFunc("/api/stores", stgHandler.PostStore).Methods(http.MethodPost)
	r.HandleFunc("/api/stores/{id}", stgHandler.PutStore).Methods(http.MethodPut)

	r.HandleFunc("/api/customers", stgHandler.GetCustomers).Methods(http.MethodGet)
	r.HandleFunc("/api/customers/{id}", stgHandler.GetCustomer).Methods(http.MethodGet)
	r.HandleFunc("/api/customers/{id}", stgHandler.DeleteCustomer).Methods(http.MethodDelete)
	r.HandleFunc("/api/customers", stgHandler.PostCustomer).Methods(http.MethodPost)
	r.HandleFunc("/api/customers/{id}", stgHandler.PutCustomer).Methods(http.MethodPut)

	return r
}

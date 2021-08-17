package server

import (
	"apigo/internal/handlers"
	"apigo/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(stgStore storage.StoreStorage, stgCustomer storage.CustomerStorage, stgCase storage.CaseStorage) http.Handler {
	r := mux.NewRouter()
	stgHandler := &handlers.StorageHandler{StgStore: stgStore, StgCustomer: stgCustomer, StgCase: stgCase}
	r.HandleFunc("/api/stores", stgHandler.GetStores).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/stores/{id}", stgHandler.GetStore).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/stores/{id}", stgHandler.DeleteStore).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/stores", stgHandler.PostStore).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/stores/{id}", stgHandler.PutStore).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/stores/{id}/cases", stgHandler.GetCasesByStoreID).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/stores/{id}/customers", stgHandler.GetCustomersByStoreID).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/api/customers", stgHandler.GetCustomers).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/customers/{id}", stgHandler.GetCustomer).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/customers/{id}", stgHandler.DeleteCustomer).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/customers", stgHandler.PostCustomer).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/customers/{id}", stgHandler.PutCustomer).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/customers/{id}/cases", stgHandler.GetCasesByCustomerID).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/customers/{id}/stores", stgHandler.GetStoreByCustomerID).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/api/cases", stgHandler.GetCases).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/cases/{id}", stgHandler.GetCase).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/cases/{id}", stgHandler.DeleteCase).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/cases", stgHandler.PostCase).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/cases/{id}", stgHandler.PutCase).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/cases/{id}/customers", stgHandler.GetCustomerByCaseID).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/cases/{id}/stores", stgHandler.GetStoreByCaseID).Methods(http.MethodGet, http.MethodOptions)
	return r
}

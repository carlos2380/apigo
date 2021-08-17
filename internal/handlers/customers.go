package handlers

import (
	"apigo/api"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (sHandler *StorageHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		retCustomers, err := sHandler.StgCustomer.GetCustomers()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Bad Request"})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(retCustomers)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		params := mux.Vars(r)
		id := params["id"]
		retCustomer, err := sHandler.StgCustomer.GetCustomer(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Item Not Found"})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(retCustomer)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		params := mux.Vars(r)
		id := params["id"]
		err := sHandler.StgCustomer.DeleteCustomer(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Item Not Found"})
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) PostCustomer(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		var customer api.Customer
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&customer)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Bad Request"})
			return
		}
		id, err := sHandler.StgCustomer.PostCustomer(&customer)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Error on update DB"})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(
			struct {
				ID string `json:"id"`
			}{ID: id})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) PutCustomer(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodPut:
		var customer api.Customer
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&customer)
		if err != nil || customer.ValidReq() != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Bad Request"})
			return
		}
		err = sHandler.StgCustomer.PutCustomer(&customer)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Error on update DB"})
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) GetCustomersByStoreID(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Bad Request"})
			return
		}

		retCases, err := sHandler.StgCustomer.GetCustomersByStoreID(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Bad Request"})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(retCases)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) GetStoreByCustomerID(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:

		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Bad Request"})
			return
		}
		retCustomer, err := sHandler.StgCustomer.GetCustomer(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Item Not Found"})
			return
		}
		params["id"] = retCustomer.StoreID
		sHandler.GetStore(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

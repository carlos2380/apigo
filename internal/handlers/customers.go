package handlers

import (
	"apigo/api"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (sHandler *StorageHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	retCustomers, err := sHandler.StgCustomer.GetCustomers()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(retCustomers) > 0 {
		json.NewEncoder(w).Encode(retCustomers)
	} else {
		json.NewEncoder(w).Encode("[]")
	}
}

func (sHandler *StorageHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	retCustomer, err := sHandler.StgCustomer.GetCustomer(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Item Not Found"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(retCustomer)

}

func (sHandler *StorageHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	err := sHandler.StgCustomer.DeleteCustomer(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Item Not Found"})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (sHandler *StorageHandler) PostCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer api.Customer
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&customer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	id, err := sHandler.StgCustomer.PostCustomer(&customer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Error on update DB"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		struct {
			Id string `json:"id"`
		}{Id: id})
}

func (sHandler *StorageHandler) PutCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer api.Customer
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&customer)
	if err != nil || customer.ValidReq() != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	err = sHandler.StgCustomer.PutCustomer(&customer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Error on update DB"})
		return
	}
	w.WriteHeader(http.StatusOK)
}

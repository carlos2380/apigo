package handlers

import (
	"apigo/api"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var stores []api.Store

func GetStores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

func GetStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	for _, item := range stores {
		if item.ID == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(
		struct {
			Error string `json:"error"`
		}{Error: "Item Not Found"})
}

func DeleteStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	for i, store := range stores {
		if store.ID == id {
			stores = append(stores[:i], stores[i+1:]...)
			w.WriteHeader(http.StatusOK)
			//json.NewEncoder(w).Encode(stores)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(
		struct {
			Error string `json:"error"`
		}{Error: "Item Not Found"})
}

func PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var store api.Store
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&store)
	if err != nil || store.ValidReq() != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	stores = append(stores, store)
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode("")

}

func PutStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

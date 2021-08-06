package handlers

import (
	"apigo/api"
	"encoding/json"
	"fmt"
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
	w.Write([]byte("Item not found"))
}

func DeleteStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	for i, store := range stores {
		if store.ID == id {
			stores = append(stores[:i], stores[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stores)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Item not found"))
}

func PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var store api.Store
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	stores = append(stores, store)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stores)
}

func PutStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

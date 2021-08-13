package handlers

import (
	"apigo/api"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (sHandler *StorageHandler) GetStores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	retStores, err := sHandler.StgStore.GetStores()
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
	_ = json.NewEncoder(w).Encode(retStores)
}

func (sHandler *StorageHandler) GetStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	retStore, err := sHandler.StgStore.GetStore(id)
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
	_ = json.NewEncoder(w).Encode(retStore)
}

func (sHandler *StorageHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	err := sHandler.StgStore.DeleteStore(id)
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
}

func (sHandler *StorageHandler) PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var store api.Store
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&store)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	id, err := sHandler.StgStore.PostStore(&store)
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
}

func (sHandler *StorageHandler) PutStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var store api.Store
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&store)
	if err != nil || store.ValidReq() != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	err = sHandler.StgStore.PutStore(&store)
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
}

package handlers

import (
	"apigo/api"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (sHandler *StorageHandler) GetCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	retCases, err := sHandler.StgCase.GetCases()
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
	if len(retCases) > 0 {
		json.NewEncoder(w).Encode(retCases)
	} else {
		json.NewEncoder(w).Encode("[]")
	}
}

func (sHandler *StorageHandler) GetCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	retCase, err := sHandler.StgCase.GetCase(id)
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
	json.NewEncoder(w).Encode(retCase)

}

func (sHandler *StorageHandler) DeleteCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	err := sHandler.StgCase.DeleteCase(id)
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

func (sHandler *StorageHandler) PostCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cas api.Case
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&cas)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	id, err := sHandler.StgCase.PostCase(&cas)
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

func (sHandler *StorageHandler) PutCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cas api.Case
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&cas)
	if err != nil || cas.ValidReq() != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	err = sHandler.StgCase.PutCase(&cas)
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

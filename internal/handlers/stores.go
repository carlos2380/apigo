package handlers

import (
	"apigo/api"
	"apigo/internal/storage"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type EnvHandler struct {
	CtrlDB storage.DataBase
}

func (env *EnvHandler) GetStores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	retStores, err := env.CtrlDB.GetStores()
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
	json.NewEncoder(w).Encode(retStores)
}

func (env *EnvHandler) GetStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	retStore, err := env.CtrlDB.GetStore(id)
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
	json.NewEncoder(w).Encode(retStore)

}

func (env *EnvHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	rowsaffected, err := env.CtrlDB.DeleteStore(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Item Not Found"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		struct {
			Rowsaffected string `json:"rows_affected"`
		}{Rowsaffected: rowsaffected})
}

func (env *EnvHandler) PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var store api.Store
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&store)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Error string `json:"error"`
			}{Error: "Bad Request"})
		return
	}
	id, err := env.CtrlDB.PostStore(&store)
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

func (env *EnvHandler) PutStore(w http.ResponseWriter, r *http.Request) {
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
	err = env.CtrlDB.PutStore(&store)
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

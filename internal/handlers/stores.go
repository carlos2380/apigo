package handlers

import (
	"apigo/api"
	"encoding/json"
	"net/http"
)

var stores []api.Store

func GetStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stores = append(stores, api.Store{ID: "1", Name: "Store1"})
	stores = append(stores, api.Store{ID: "2", Name: "Store2"})
	json.NewEncoder(w).Encode(stores)
}

package handlers

import (
	"apigo/api"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (sHandler *StorageHandler) GetCases(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		retCases, err := sHandler.StgCase.GetCases()
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

func (sHandler *StorageHandler) GetCase(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:

		params := mux.Vars(r)
		id := params["id"]
		retCase, err := sHandler.StgCase.GetCase(id)
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
		_ = json.NewEncoder(w).Encode(retCase)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) DeleteCase(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		params := mux.Vars(r)
		id := params["id"]

		err := sHandler.StgCase.DeleteCase(id)
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

func (sHandler *StorageHandler) PostCase(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		var cas api.Case
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&cas)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Bad Request"})
			return
		}
		id, err := sHandler.StgCase.PostCase(&cas)
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

func (sHandler *StorageHandler) PutCase(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodPut:
		var cas api.Case
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&cas)
		if err != nil || cas.ValidReq() != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Bad Request"})
			return
		}
		err = sHandler.StgCase.PutCase(&cas)
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

func (sHandler *StorageHandler) GetCasesByStoreID(w http.ResponseWriter, r *http.Request) {
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

		retCases, err := sHandler.StgCase.GetCasesByStoreID(id)
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
		_ = json.NewEncoder(w).Encode(&api.CasesJSON{Cases: retCases})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) GetCasesByCustomerID(w http.ResponseWriter, r *http.Request) {
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

		retCases, err := sHandler.StgCase.GetCasesByCustomerID(id)
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

		_ = json.NewEncoder(w).Encode(&api.CasesJSON{Cases: retCases})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) GetCustomerByCaseID(w http.ResponseWriter, r *http.Request) {
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
		retCase, err := sHandler.StgCase.GetCase(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Item Not Found"})
			return
		}
		params["id"] = retCase.CustomerID
		sHandler.GetCustomer(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (sHandler *StorageHandler) GetStoreByCaseID(w http.ResponseWriter, r *http.Request) {
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
		retCase, err := sHandler.StgCase.GetCase(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(
				struct {
					Error string `json:"error"`
				}{Error: "Item Not Found"})
			return
		}
		params["id"] = retCase.StoreID
		sHandler.GetStore(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

package handlers

import "apigo/internal/storage"

type StorageHandler struct {
	StgStore    storage.StoreStorage
	StgCustomer storage.CustomerStorage
	StgCase     storage.CaseStorage
}

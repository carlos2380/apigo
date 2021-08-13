package storage

import (
	"apigo/api"
)

type Store struct{}
type StoreStorage interface {
	GetStores() (*api.StoresJSON, error)
	GetStore(storeID string) (*api.Store, error)
	DeleteStore(storeID string) error
	PostStore(storeReq *api.Store) (string, error)
	PutStore(storeReq *api.Store) error
}

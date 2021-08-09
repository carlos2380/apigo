package storage

import "apigo/api"

type Store struct{}
type DataBase interface {
	InitDB() error
	CloseB() error
	GetStores() ([]api.Store, error)
	GetStore(storeID string) (*api.Store, error)
	DeleteStore(storeID string) error
	PostStore(storeReq *api.Store) error
	PutStore(storeReq *api.Store) error
}

/*
type Customer struct{}
type CustomerStorage interface {
	CreateCustomer() error
	ListCustomers(customerID string) ([]Customer, error)
	GetCustomer() (*Customer, error)
	DeleteCustomer() error
}

type Store struct{}
type StoreStorage interface {
	//CreateStore() error
	//ListStores(storeID string) ([]Store, error)
	GetStore() (*Store, error)
	//DeleteStore() error
}

type Case struct{}
type CaseStorage interface {
	CreateCase() error
	ListCases(caseID string) ([]Case, error)
	GetCase() (*Case, error)
	DeleteCase() error
}
*/

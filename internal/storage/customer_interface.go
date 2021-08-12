package storage

import "apigo/api"

type Customer struct{}
type CustomerStorage interface {
	GetCustomers() ([]*api.Customer, error)
	GetCustomer(customerID string) (*api.Customer, error)
	DeleteCustomer(customerID string) error
	PostCustomer(customerReq *api.Customer) (string, error)
	PutCustomer(customerReq *api.Customer) error
}

/*
type Customer struct{}
type CustomerStorage interface {
	CretateCustomer() error
	LisCustomers(customerID string) ([]Customer, error)
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

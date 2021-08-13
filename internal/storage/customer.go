package storage

import "apigo/api"

type Customer struct{}
type CustomerStorage interface {
	GetCustomers() (*api.CustomersJSON, error)
	GetCustomer(customerID string) (*api.Customer, error)
	DeleteCustomer(customerID string) error
	PostCustomer(customerReq *api.Customer) (string, error)
	PutCustomer(customerReq *api.Customer) error
}

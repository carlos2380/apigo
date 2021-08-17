package api

import "errors"

type Customer struct {
	ID        string  `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Age       *string `json:"age"`
	Email     *string `json:"email"`
	StoreID   string  `json:"store_id"`
}

type CustomersJSON struct {
	Customers []*Customer `json:"customers"`
}

func (customer *Customer) ValidReq() error {
	if customer.ID != "" {
		return nil
	}
	return errors.New("bad parameters")
}

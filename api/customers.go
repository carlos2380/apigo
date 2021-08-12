package api

import "errors"

type Customer struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       string `json:"age"`
	Email     string `json:"email"`
}

func (customer *Customer) ValidReq() error {

	if customer.Id != "" {
		return nil
	}

	return errors.New("bad parameters")
}

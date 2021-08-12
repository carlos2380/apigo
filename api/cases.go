package api

import "errors"

type Case struct {
	Id         string `json:"id"`
	StartTime  string `json:"start_time_stamp"`
	EndTime    string `json:"end_time_stamp"`
	CustomerId string `json:"customer_id"`
	StoreId    string `json:"store_id"`
}

func (cas *Case) ValidReq() error {

	if cas.Id != "" {
		return nil
	}

	return errors.New("bad parameters")
}

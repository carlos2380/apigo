package api

import "errors"

type Case struct {
	ID         string  `json:"id"`
	StartTime  *string `json:"start_time_stamp"`
	EndTime    *string `json:"end_time_stamp"`
	CustomerID string  `json:"customer_id"`
	StoreID    string  `json:"store_id"`
}

type CasesJSON struct {
	Cases []*Case `json:"cases"`
}

func (cas *Case) ValidReq() error {
	if cas.ID != "" {
		return nil
	}
	return errors.New("bad parameters")
}

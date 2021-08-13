package api

import "errors"

type Store struct {
	ID      string  `json:"id"`
	Name    *string `json:"name"`
	Address *string `json:"address"`
}

type StoresJSON struct {
	Stores []*Store `json:"stores"`
}

func (s *Store) ValidReq() error {
	if s.ID != "" {
		return nil
	}
	return errors.New("bad parameters")
}

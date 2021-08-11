package api

import "errors"

type Store struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (s *Store) ValidReq() error {

	if s.Id != "" {
		return nil
	}

	return errors.New("bad parameters")
}

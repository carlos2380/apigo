package api

import "errors"

type Store struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Store) ValidReq() error {

	if s.ID != "" && s.Name != "" {
		return nil
	}

	return errors.New("bad parameters")
}

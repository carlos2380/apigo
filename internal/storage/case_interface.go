package storage

import "apigo/api"

type Case struct{}
type CaseStorage interface {
	GetCases() ([]*api.Case, error)
	GetCase(caseID string) (*api.Case, error)
	DeleteCase(caseID string) error
	PostCase(caseReq *api.Case) (string, error)
	PutCase(caseReq *api.Case) error
	GetCasesByStoreId(storeID string) ([]*api.Case, error)
}

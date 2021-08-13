package storage

import "apigo/api"

type Case struct{}
type CaseStorage interface {
	GetCases() (*api.CasesJSON, error)
	GetCase(caseID string) (*api.Case, error)
	DeleteCase(caseID string) error
	PostCase(caseReq *api.Case) (string, error)
	PutCase(caseReq *api.Case) error
	GetCasesByStoreID(storeID string) ([]*api.Case, error)
	GetCasesByCustomerID(storeID string) ([]*api.Case, error)
}

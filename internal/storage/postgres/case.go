package postgres

import (
	"apigo/api"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type CaseDB struct {
	DB *sql.DB
}

type Case struct {
	id         int64
	startTime  *time.Time
	endTime    *time.Time
	customerID int64
	storeID    int64
	created    time.Time
	modified   time.Time
}

func (pdb *CaseDB) CloseDB() error {
	return pdb.DB.Close()
}

func (pdb *CaseDB) GetCases() (*api.CasesJSON, error) {
	rows, err := pdb.DB.Query("Select * From cases")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	cases := []*api.Case{}
	for rows.Next() {
		var c Case
		err = rows.Scan(&c.id, &c.startTime, &c.endTime, &c.customerID, &c.storeID, &c.created, &c.modified)
		if err != nil {
			return nil, err
		}
		starttime := ""
		if c.startTime != nil {
			starttime = c.startTime.String()
		}
		endtime := ""
		if c.endTime != nil {
			endtime = c.startTime.String()
		}
		cases = append(
			cases,
			&api.Case{
				ID:        strconv.FormatInt(c.id, 10),
				StartTime: &starttime, EndTime: &endtime,
				CustomerID: strconv.FormatInt(c.customerID, 10),
				StoreID:    strconv.FormatInt(c.storeID, 10),
			},
		)
	}

	return &api.CasesJSON{Cases: cases}, nil
}

func (pdb *CaseDB) GetCase(caseID string) (*api.Case, error) {
	row := pdb.DB.QueryRow(fmt.Sprintf("Select * From cases WHERE id = '%s'", caseID))
	var c Case
	err := row.Scan(&c.id, &c.startTime, &c.endTime, &c.customerID, &c.storeID, &c.created, &c.modified)
	if err != nil {
		return nil, err
	}
	starttime := ""
	if c.startTime != nil {
		starttime = c.startTime.String()
	}
	endtime := ""
	if c.endTime != nil {
		endtime = c.endTime.String()
	}
	return &api.Case{
		ID:        strconv.FormatInt(c.id, 10),
		StartTime: &starttime, EndTime: &endtime,
		CustomerID: strconv.FormatInt(c.customerID, 10),
		StoreID:    strconv.FormatInt(c.storeID, 10),
	}, nil
}

func (pdb *CaseDB) DeleteCase(caseID string) (err error) {
	sqlStatement := `
		DELETE FROM cases
		WHERE id = $1
	`
	ctx := context.Background()
	tx, err := pdb.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	resultsql, err := tx.ExecContext(ctx, sqlStatement, caseID)
	if err != nil {
		return err
	}
	rowsAffected, err := resultsql.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		err = errors.New("item no found")
		return
	}

	return nil
}

func (pdb *CaseDB) PostCase(caseReq *api.Case) (_ string, err error) {
	timeStr := time.Now().Format(time.RFC3339)
	sqlStatement := fmt.Sprintf(`INSERT INTO cases (start_time_stamp, end_time_stamp, customer_id, store_id, created_on, modified_on)
		VALUES ( '%s', '%s', '%s', '%s', '%s', '%s')
		RETURNING id`,
		*caseReq.StartTime, *caseReq.EndTime, caseReq.CustomerID, caseReq.StoreID, timeStr, timeStr)
	lastInsertid := int64(0)

	ctx := context.Background()
	tx, err := pdb.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	err = tx.QueryRowContext(ctx, sqlStatement).Scan(&lastInsertid)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(lastInsertid, 10), nil
}

func (pdb *CaseDB) PutCase(caseReq *api.Case) (err error) {
	sqlStatement := `
		UPDATE cases
		SET start_time_stamp = $1, end_time_stamp = $2, customer_id = $3, store_id = $4, modified_on = $5
		WHERE id = $6
	`
	timeStr := time.Now().Format(time.RFC3339)

	ctx := context.Background()
	tx, err := pdb.DB.BeginTx(ctx, nil)

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, sqlStatement, caseReq.StartTime, caseReq.EndTime, caseReq.CustomerID, caseReq.StoreID, timeStr, caseReq.ID)
	if err != nil {
		return err
	}

	return nil
}

func (pdb *CaseDB) GetCasesByStoreID(storeID string) ([]*api.Case, error) {
	rows, err := pdb.DB.Query(fmt.Sprintf("Select * From cases Where store_id = '%s'", storeID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	cases := []*api.Case{}
	for rows.Next() {
		var c Case
		err = rows.Scan(&c.id, &c.startTime, &c.endTime, &c.customerID, &c.storeID, &c.created, &c.modified)
		if err != nil {
			return nil, err
		}

		starttime := ""
		if c.startTime != nil {
			starttime = c.startTime.String()
		}
		endtime := ""
		if c.endTime != nil {
			endtime = c.endTime.String()
		}
		cases = append(
			cases,
			&api.Case{
				ID:         strconv.FormatInt(c.id, 10),
				StartTime:  &starttime,
				EndTime:    &endtime,
				CustomerID: strconv.FormatInt(c.customerID, 10),
				StoreID:    strconv.FormatInt(c.storeID, 10),
			},
		)
	}
	return cases, nil
}

func (pdb *CaseDB) GetCasesByCustomerID(customerID string) ([]*api.Case, error) {
	rows, err := pdb.DB.Query(fmt.Sprintf("Select * From cases Where customer_id = '%s'", customerID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	var cases []*api.Case
	for rows.Next() {
		var c Case
		err = rows.Scan(&c.id, &c.startTime, &c.endTime, &c.customerID, &c.storeID, &c.created, &c.modified)
		if err != nil {
			return nil, err
		}

		starttime := ""
		if c.startTime != nil {
			starttime = c.startTime.String()
		}
		endtime := ""
		if c.endTime != nil {
			endtime = c.endTime.String()
		}
		cases = append(
			cases,
			&api.Case{
				ID:         strconv.FormatInt(c.id, 10),
				StartTime:  &starttime,
				EndTime:    &endtime,
				CustomerID: strconv.FormatInt(c.customerID, 10),
				StoreID:    strconv.FormatInt(c.storeID, 10),
			},
		)
	}
	return cases, nil
}

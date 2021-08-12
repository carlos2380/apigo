package postgres

import (
	"apigo/api"
	"apigo/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type CaseDB struct {
	Db *sql.DB
}

type Case struct {
	id         int64
	startTime  time.Time
	endTime    time.Time
	customerID int64
	storeID    int64
	created    time.Time
	modified   time.Time
}

func NewPostgresCase() (storage.CaseStorage, error) {
	uri := "user=postgres password=secret host=172.17.0.1 port=5432 dbname=postgres connect_timeout=20 sslmode=disable"
	var err error
	postgresDB, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	if err := postgresDB.Ping(); err != nil {
		return nil, err

	}
	return &CaseDB{Db: postgresDB}, nil
}

func (pdb *CaseDB) CloseDB() error {
	return pdb.Db.Close()
}

func (pdb *CaseDB) GetCases() ([]*api.Case, error) {
	rows, err := pdb.Db.Query("Select * From kase")
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
		cases = append(cases, &api.Case{Id: strconv.FormatInt(c.id, 10), StartTime: c.startTime.String(), EndTime: c.endTime.String(), CustomerId: strconv.FormatInt(c.customerID, 10), StoreId: strconv.FormatInt(c.storeID, 10)})
	}
	return cases, nil
}

func (pdb *CaseDB) GetCase(caseID string) (*api.Case, error) {
	row := pdb.Db.QueryRow("Select * From kase WHERE id = '" + caseID + "'")
	var c Case
	err := row.Scan(&c.id, &c.startTime, &c.endTime, &c.customerID, &c.storeID, &c.created, &c.modified)
	if err != nil {
		return nil, err
	}
	return &api.Case{Id: strconv.FormatInt(c.id, 10), StartTime: c.startTime.String(), EndTime: c.endTime.String(), CustomerId: strconv.FormatInt(c.customerID, 10), StoreId: strconv.FormatInt(c.storeID, 10)}, nil
}

func (pdb *CaseDB) DeleteCase(caseID string) (err error) {
	sqlStatement := `
		DELETE FROM kase
		WHERE id = $1
	`
	ctx := context.Background()
	tx, err := pdb.Db.BeginTx(ctx, nil)
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
		tx.Rollback()
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
	sqlStatement := fmt.Sprintf(`INSERT INTO kase (start_time_stamp, end_time_stamp, customer_id, store_id, created_on, modified_on)
		VALUES ( '%s', '%s', '%s', '%s', '%s', '%s')
		RETURNING id`,
		caseReq.StartTime, caseReq.EndTime, caseReq.CustomerId, caseReq.StoreId, timeStr, timeStr)
	lastInsertid := int64(0)

	ctx := context.Background()
	tx, err := pdb.Db.BeginTx(ctx, nil)
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
		UPDATE kase
		SET start_time_stamp = $1, end_time_stamp = $2, customer_id = $3, store_id = $4, modified_on = $5
		WHERE id = $6
	`
	timeStr := time.Now().Format(time.RFC3339)

	ctx := context.Background()
	tx, err := pdb.Db.BeginTx(ctx, nil)

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
	_, err = tx.ExecContext(ctx, sqlStatement, caseReq.StartTime, caseReq.EndTime, caseReq.CustomerId, caseReq.StoreId, timeStr, caseReq.Id)
	if err != nil {
		return err
	}

	return nil
}

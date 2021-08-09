package handlers

import (
	"apigo/internal/storage/postgres"
)

type EnvHandler struct {
	CtrlDB *postgres.PostgresDB
}

/*func (env *EnvHandler) InitDB() error {
	var err error
	env.db, err = postgres.DataBase.InitDB()
	return err
}

func (hdb *HandlerDB) CloseDB() error {
	hdb.db.CloseDB()
	return nil
}*/

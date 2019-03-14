package dbr

import "database/sql"

type eventHandler func(operation SqlOperation, table []string, query string, err error, rows *sql.Rows, sqlResult sql.Result)

type SuccessEventHandler func(operation SqlOperation, table []string, query string, rows *sql.Rows, sqlResult sql.Result)
type ErrorEventHandler func(operation SqlOperation, table []string, query string, err error)

func NewDb(database database, dialect dialect) *db {
	return &db{
		database: database,
		Dialect:  dialect,
	}
}

func (dbr *Dbr) handle(operation SqlOperation, table []string, query string, err error, rows *sql.Rows, sqlResult sql.Result) {
	if err == nil && dbr.isSuccessEventHandlerActive {
		dbr.successEventHandler(operation, table, query, rows, sqlResult)
	}

	if err != nil && dbr.isErrorEventHandlerActive {
		dbr.errorEventHandler(operation, table, query, err)
	}
}

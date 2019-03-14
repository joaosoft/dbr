package dbr

type eventHandler func(operation SqlOperation, table []string, query string, err error, result interface{})

type SuccessEventHandler func(operation SqlOperation, table []string, query string, result interface{})
type ErrorEventHandler func(operation SqlOperation, table []string, query string, err error)

func NewDb(database database, dialect dialect) *db {
	return &db{
		database: database,
		Dialect:  dialect,
	}
}

func (dbr *Dbr) handle(operation SqlOperation, table []string, query string, err error, result interface{}) {
	if err == nil && dbr.isSuccessEventHandlerActive {
		dbr.successEventHandler(operation, table, query, result)
	}

	if err != nil && dbr.isErrorEventHandlerActive {
		dbr.errorEventHandler(operation, table, query, err)
	}
}

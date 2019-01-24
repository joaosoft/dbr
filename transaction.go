package dbr

import (
	"database/sql"
)

type Transaction struct {
	commited bool
	dbr      *Dbr
	db       *db
}

func newTransaction(dbr *Dbr, db *db) *Transaction {
	return &Transaction{dbr: dbr, db: db}
}

func (tx *Transaction) Commit() error {
	err := tx.db.database.(*sql.Tx).Commit()
	if err != nil {
		return err
	}

	tx.commited = true
	return nil
}

func (tx *Transaction) Rollback() error {
	return tx.db.database.(*sql.Tx).Rollback()
}

func (tx *Transaction) RollbackUnlessCommit() error {
	if !tx.commited {
		return tx.db.database.(*sql.Tx).Rollback()
	}

	return nil
}

func (tx *Transaction) Select(column ...string) *StmtSelect {
	return newStmtSelect(tx.dbr, tx.dbr.connections.write, &StmtWith{}, column)
}

func (tx *Transaction) Insert() *StmtInsert {
	return newStmtInsert(tx.dbr, tx.dbr.connections.write, &StmtWith{})
}

func (tx *Transaction) Update(table string) *StmtUpdate {
	return newStmtUpdate(tx.dbr, tx.dbr.connections.write, &StmtWith{}, table)
}

func (tx *Transaction) Delete() *StmtDelete {
	return newStmtDelete(tx.dbr, tx.dbr.connections.write, &StmtWith{})
}

func (tx *Transaction) Execute(query string) *StmtExecute {
	return newStmtExecute(tx.dbr, tx.dbr.connections.write, query)
}

func (tx *Transaction) With(name string, builder builder) *StmtWith {
	return newStmtWith(tx.dbr, tx.dbr.connections, name, false, builder)
}

func (tx *Transaction) WithRecursive(name string, builder builder) *StmtWith {
	return newStmtWith(tx.dbr, tx.dbr.connections, name, true, builder)
}

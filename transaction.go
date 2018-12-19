package dbr

import (
	"database/sql"
)

type Transaction struct {
	db       *db
	commited bool
}

func newTransaction(db *db) *Transaction {
	return &Transaction{db: db}
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
	return newStmtSelect(tx.db, &StmtWith{}, column)
}

func (tx *Transaction) Insert() *StmtInsert {
	return newStmtInsert(tx.db, &StmtWith{})
}

func (tx *Transaction) Update(table string) *StmtUpdate {
	return newStmtUpdate(tx.db, &StmtWith{}, table)
}

func (tx *Transaction) Delete() *StmtDelete {
	return newStmtDelete(tx.db, &StmtWith{})
}

func (tx *Transaction) Execute(query string) *StmtExecute {
	return newStmtExecute(tx.db, query)
}

func (tx *Transaction) With(name string, builder builder) *StmtWith {
	return newStmtWith(&connections{read: tx.db, write: tx.db}, name, false, builder)
}

func (tx *Transaction) WithRecursive(name string, builder builder) *StmtWith {
	return newStmtWith(&connections{read: tx.db, write: tx.db}, name, true, builder)
}

package dbr

import (
	"database/sql"
)

type Transaction struct {
	db       *Db
	commited bool
}

func newTransaction(db *Db) *Transaction {
	return &Transaction{db: db}
}

func (tx *Transaction) Commit() error {
	err := tx.db.Database.(*sql.Tx).Commit()
	if err != nil {
		return err
	}

	tx.commited = true
	return nil
}

func (tx *Transaction) Rollback() error {
	return tx.db.Database.(*sql.Tx).Rollback()
}

func (tx *Transaction) RollbackUnlessCommit() error {
	if !tx.commited {
		return tx.db.Database.(*sql.Tx).Rollback()
	}

	return nil
}

func (tx *Transaction) Select(column ...string) *StmtSelect {
	return newStmtSelect(tx.db, column)
}

func (tx *Transaction) Insert() *StmtInsert {
	return newStmtInsert(tx.db)
}

func (tx *Transaction) Update(table string) *StmtUpdate {
	return newStmtUpdate(tx.db, table)
}

func (tx *Transaction) Delete() *StmtDelete {
	return newStmtDelete(tx.db)
}

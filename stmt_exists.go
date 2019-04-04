package dbr

import (
	"fmt"
)

type StmtExists struct {
	stmtSelect *StmtSelect
	isNot      bool

	db *db
}

func newStmtExists(db *db, stmtSelect *StmtSelect, isNot bool) *StmtExists {
	return &StmtExists{
		db:         db,
		stmtSelect: stmtSelect,
		isNot:      isNot,
	}
}

func (stmt *StmtExists) Build() (string, error) {

	var query string

	stmtSelect, err := stmt.stmtSelect.Build()
	if err != nil {
		return "", err
	}

	if stmt.isNot {
		query += fmt.Sprintf("%s ", constFunctionNot)
	}

	query += fmt.Sprintf("%s %s", constFunctionExists, stmtSelect)

	return query, nil
}

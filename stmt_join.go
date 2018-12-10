package dbr

import (
	"fmt"
)

type StmtJoin struct {
	join  Join
	table string
	on    string

	db *Db
}

type Join string

const (
	ConstJoin      Join = "JOIN"
	ConstLeftJoin  Join = "LEFT JOIN"
	ConstRightJoin Join = "RIGHT JOIN"
	ConstFullJoin  Join = "FULL JOIN"
)

func newStmtJoin(db *Db, join Join, table, on string) *StmtJoin {
	return &StmtJoin{db: db, table: table, on: on}
}

func (stmt *StmtJoin) Build() (string, error) {
	query := fmt.Sprintf("%s %s ON (%s)", stmt.join, stmt.table, stmt.on)

	return query, nil
}

package dbr

import (
	"fmt"
	"strings"
)

type onConflictType string
type onConflictDoType string

const (
	onConflictColumn     onConflictType = "column"
	onConflictConstraint onConflictType = "constraint"

	onConflictDoNothing onConflictDoType = "nothing"
	onConflictDoUpdate  onConflictDoType = "update"
)

type StmtConflict struct {
	table              string
	onConflictType     onConflictType
	onConflict         columns
	onConflictDoType   onConflictDoType
	onConflictDoUpdate sets

	db *db
}

func newStmtConflict(db *db, table string) *StmtConflict {
	return &StmtConflict{db: db, table: table}
}

func (stmt *StmtConflict) Build() (string, error) {

	if stmt.onConflictType == "" {
		return "", nil
	}

	query := " ON CONFLICT "
	switch stmt.onConflictType {
	case onConflictColumn:
		query += fmt.Sprintf("(%s) ", strings.Join(stmt.onConflict, ", "))
	case onConflictConstraint:
		query += fmt.Sprintf("ON CONSTRAINT (%s) ", strings.Join(stmt.onConflict, ", "))
	}

	switch stmt.onConflictDoType {
	case onConflictDoNothing:
		query += "DO NOTHING"
	case onConflictDoUpdate:
		sets, err := stmt.onConflictDoUpdate.Build()
		if err != nil {
			return "", err
		}
		query += fmt.Sprintf("DO UPDATE SET %s", sets)
	}

	return query, nil
}

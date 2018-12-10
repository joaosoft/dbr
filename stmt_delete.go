package dbr

import (
	"database/sql"
	"fmt"
)

type StmtDelete struct {
	table      string
	conditions conditions
	returning  columns

	db *Db
}

func newStmtDelete(db *Db) *StmtDelete {
	return &StmtDelete{db: db, conditions: conditions{db: db}}
}

func (stmt *StmtDelete) From(table string) *StmtDelete {
	stmt.table = table
	return stmt
}

func (stmt *StmtDelete) Where(query string, valuesList ...interface{}) *StmtDelete {
	stmt.conditions.list = append(stmt.conditions.list, &condition{query: query, values: values{list: valuesList, db: stmt.db}})
	return stmt
}

func (stmt *StmtDelete) Build() (string, error) {
	query := fmt.Sprintf("DELETE FROM %s", stmt.table)

	if len(stmt.conditions.list) > 0 {
		conds, err := stmt.conditions.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" WHERE %s", conds)
	}

	if len(stmt.returning) > 0 {
		returning, err := stmt.returning.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" RETURNING %s", returning)
	}

	return query, nil
}

func (stmt *StmtDelete) Exec() (sql.Result, error) {
	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	return stmt.db.Exec(query)
}

func (stmt *StmtDelete) Return(column ...string) *StmtDelete {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

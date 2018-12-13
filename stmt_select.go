package dbr

import (
	"fmt"
	"reflect"
)

type StmtSelect struct {
	columns           columns
	tables            tables
	joins             joins
	conditions        conditions
	isDistinct        bool
	distinctColumns   columns
	distinctOnColumns columns
	returning         columns

	db *Db
}

func newStmtSelect(db *Db, columns []string) *StmtSelect {
	return &StmtSelect{db: db, columns: columns, conditions: conditions{db: db}}
}

func (stmt *StmtSelect) From(tables ...string) *StmtSelect {
	stmt.tables = tables
	return stmt
}

func (stmt *StmtSelect) Where(query string, values ...interface{}) *StmtSelect {
	stmt.conditions.list = append(stmt.conditions.list, &condition{query: query, values: values})
	return stmt
}

func (stmt *StmtSelect) Join(table, on string) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.db, ConstJoin, table, on))
	return stmt
}

func (stmt *StmtSelect) LeftJoin(table, on string) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.db, ConstLeftJoin, table, on))
	return stmt
}

func (stmt *StmtSelect) RightJoin(table, on string) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.db, ConstRightJoin, table, on))
	return stmt
}

func (stmt *StmtSelect) FullJoin(table, on string) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.db, ConstFullJoin, table, on))
	return stmt
}

func (stmt *StmtSelect) Distinct(column ...string) *StmtSelect {
	stmt.isDistinct = true
	stmt.distinctColumns = append(stmt.distinctColumns, column...)
	return stmt
}

func (stmt *StmtSelect) DistinctOn(column ...string) *StmtSelect {
	stmt.distinctOnColumns = append(stmt.distinctOnColumns, column...)
	return stmt
}

func (stmt *StmtSelect) Return(column ...string) *StmtSelect {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtSelect) Build() (string, error) {

	var distinct string
	if stmt.isDistinct {
		distinct = "DISTINCT "
	}

	distinctColumns, err := stmt.distinctColumns.Build()
	if err != nil {
		return "", err
	}

	var distinctOn string
	if stmt.isDistinct {
		distinctOn = "DISTINCT ON (%s) "
	}

	distinctOnColumns, err := stmt.distinctOnColumns.Build()
	if err != nil {
		return "", err
	}

	columns, err := stmt.columns.Build()
	if err != nil {
		return "", err
	}

	tables, err := stmt.tables.Build()
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT %s%s%s%s%s FROM %s", distinct, distinctColumns, distinctOn, distinctOnColumns, columns, tables)

	if len(stmt.joins) > 0 {
		joins, err := stmt.joins.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" %s", joins)
	}

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

func (stmt *StmtSelect) Load(object interface{}) (int, error) {
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return 0, ErrorInvalidPointer
	}

	query, err := stmt.Build()
	if err != nil {
		return 0, err
	}

	rows, err := stmt.db.Query(query)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	return read(columns, rows, value)
}

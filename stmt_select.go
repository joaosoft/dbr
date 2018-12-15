package dbr

import (
	"fmt"
	"reflect"
)

type StmtSelect struct {
	withs             withs
	columns           columns
	tables            tables
	joins             joins
	conditions        conditions
	isDistinct        bool
	distinctColumns   columns
	distinctOnColumns columns
	orders            orders
	unions            unions
	returning         columns

	db *db
}

func newStmtSelect(db *db, withs withs, columns []string) *StmtSelect {
	return &StmtSelect{db: db, withs: withs, columns: columns, conditions: conditions{db: db}}
}

func (stmt *StmtSelect) From(tables ...string) *StmtSelect {
	stmt.tables = append(stmt.tables, tables...)
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

func (stmt *StmtSelect) Union(stmtUnion *StmtSelect) *StmtSelect {
	stmt.unions = append(stmt.unions, stmtUnion)
	return stmt
}

func (stmt *StmtSelect) OrderAsc(columns ...string) *StmtSelect {
	for _, column := range columns {
		stmt.orders = append(stmt.orders, &order{column: column, direction: orderAsc})
	}

	return stmt
}

func (stmt *StmtSelect) OrderDesc(columns ...string) *StmtSelect {
	for _, column := range columns {
		stmt.orders = append(stmt.orders, &order{column: column, direction: orderDesc})
	}

	return stmt
}

func (stmt *StmtSelect) Return(column ...string) *StmtSelect {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtSelect) Build() (string, error) {
	var query string

	// withs
	if len(stmt.withs) > 0 {
		withs, err := stmt.withs.Build()
		if err != nil {
			return "", err
		}
		query += fmt.Sprintf("WITH %s ", withs)
	}

	// distinct
	var distinct string
	if stmt.isDistinct {
		distinct = "DISTINCT "
	}

	distinctColumns, err := stmt.distinctColumns.Build()
	if err != nil {
		return "", err
	}

	// distinct on
	var distinctOn string
	if stmt.isDistinct {
		distinctOn = "DISTINCT ON (%s) "
	}

	distinctOnColumns, err := stmt.distinctOnColumns.Build()
	if err != nil {
		return "", err
	}

	// columns
	columns, err := stmt.columns.Build()
	if err != nil {
		return "", err
	}

	// tables
	tables, err := stmt.tables.Build()
	if err != nil {
		return "", err
	}

	// query
	query += fmt.Sprintf("SELECT %s%s%s%s%s FROM %s", distinct, distinctColumns, distinctOn, distinctOnColumns, columns, tables)

	// joins
	if len(stmt.joins) > 0 {
		joins, err := stmt.joins.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" %s", joins)
	}

	// conditions
	if len(stmt.conditions.list) > 0 {
		conds, err := stmt.conditions.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" WHERE %s", conds)
	}

	// unions
	if len(stmt.unions) > 0 {
		unions, err := stmt.unions.Build()
		if err != nil {
			return "", err
		}

		query += unions
	}

	// orders
	if len(stmt.orders) > 0 {
		orders, err := stmt.orders.Build()
		if err != nil {
			return "", err
		}

		query += orders
	}

	// returning
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

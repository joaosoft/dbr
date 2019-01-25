package dbr

import (
	"fmt"
	"reflect"
	"time"
)

type StmtSelect struct {
	withStmt          *StmtWith
	columns           columns
	tables            tables
	joins             joins
	conditions        conditions
	isDistinct        bool
	distinctColumns   columns
	distinctOnColumns columns
	unions            unions
	groupBy           groupBy
	having            conditions
	orders            orders
	returning         columns
	limit             int
	offset            int

	Dbr      *Dbr
	Db       *db
	Duration time.Duration
}

func newStmtSelect(dbr *Dbr, db *db, withStmt *StmtWith, columns []interface{}) *StmtSelect {
	return &StmtSelect{Dbr: dbr, Db: db, withStmt: withStmt, columns: columns, conditions: conditions{db: dbr.Connections.Read}, having: conditions{db: dbr.Connections.Read}}
}

func (stmt *StmtSelect) From(tables ...interface{}) *StmtSelect {
	stmt.tables = append(stmt.tables, tables...)
	return stmt
}

func (stmt *StmtSelect) Where(query string, values ...interface{}) *StmtSelect {
	stmt.conditions.list = append(stmt.conditions.list, &condition{query: query, values: values})
	return stmt
}

func (stmt *StmtSelect) Join(table, on string) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, ConstJoin, table, on))
	return stmt
}

func (stmt *StmtSelect) LeftJoin(table, on string) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, ConstLeftJoin, table, on))
	return stmt
}

func (stmt *StmtSelect) RightJoin(table, on string) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, ConstRightJoin, table, on))
	return stmt
}

func (stmt *StmtSelect) FullJoin(table, on string) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, ConstFullJoin, table, on))
	return stmt
}

func (stmt *StmtSelect) Distinct(column ...interface{}) *StmtSelect {
	stmt.isDistinct = true
	stmt.distinctColumns = append(stmt.distinctColumns, column...)
	return stmt
}

func (stmt *StmtSelect) DistinctOn(column ...interface{}) *StmtSelect {
	stmt.distinctOnColumns = append(stmt.distinctOnColumns, column...)
	return stmt
}

func (stmt *StmtSelect) Union(stmtUnion *StmtSelect) *StmtSelect {
	stmt.unions = append(stmt.unions, &union{unionType: unionNormal, stmt: stmtUnion})
	return stmt
}

func (stmt *StmtSelect) Intersect(stmtUnion *StmtSelect) *StmtSelect {
	stmt.unions = append(stmt.unions, &union{unionType: unionIntersect, stmt: stmtUnion})
	return stmt
}

func (stmt *StmtSelect) Except(stmtUnion *StmtSelect) *StmtSelect {
	stmt.unions = append(stmt.unions, &union{unionType: unionExcept, stmt: stmtUnion})
	return stmt
}

func (stmt *StmtSelect) GroupBy(columns ...string) *StmtSelect {
	stmt.groupBy = append(stmt.groupBy, columns...)
	return stmt
}

func (stmt *StmtSelect) Having(query string, values ...interface{}) *StmtSelect {
	stmt.having.list = append(stmt.having.list, &condition{query: query, values: values})
	return stmt
}

func (stmt *StmtSelect) OrderBy(column string, direction direction) *StmtSelect {
	stmt.orders = append(stmt.orders, &order{column: column, direction: direction})
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

func (stmt *StmtSelect) Return(column ...interface{}) *StmtSelect {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtSelect) Limit(limit int) *StmtSelect {
	stmt.limit = limit
	return stmt
}

func (stmt *StmtSelect) Offset(offset int) *StmtSelect {
	stmt.offset = offset
	return stmt
}

func (stmt *StmtSelect) Build() (string, error) {
	var query string

	// withStmt
	if len(stmt.withStmt.withs) > 0 {
		withStmt, err := stmt.withStmt.Build()
		if err != nil {
			return "", err
		}
		query += withStmt
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

	// group by
	if len(stmt.groupBy) > 0 {
		groupBy, err := stmt.groupBy.Build()
		if err != nil {
			return "", err
		}

		query += groupBy
	}

	// having
	if len(stmt.having.list) > 0 {
		havingConds, err := stmt.having.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" HAVING %s", havingConds)
	}

	// orders
	if len(stmt.orders) > 0 {
		orders, err := stmt.orders.Build()
		if err != nil {
			return "", err
		}

		query += orders
	}

	// limit
	if stmt.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", stmt.limit)
	}

	// offset
	if stmt.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", stmt.offset)
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

	if reflect.ValueOf(object).Kind() != reflect.Ptr {
		panic("the object is not a pointer the load")
	}

	startTime := time.Now()
	defer func() {
		stmt.Duration = time.Since(startTime)
	}()

	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return 0, ErrorInvalidPointer
	}

	query, err := stmt.Build()
	if err != nil {
		return 0, err
	}

	rows, err := stmt.Db.Query(query)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	var cols []interface{}
	for _, col := range columns {
		cols = append(cols, col)
	}

	return read(cols, rows, value)
}

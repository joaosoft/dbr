package dbr

import (
	"fmt"
)

type StmtWith struct {
	withs withs

	conn *connections
}

func newStmtWith(conn *connections, name string, builder builder) *StmtWith {
	return &StmtWith{conn: conn, withs: withs{&with{name: name, builder: builder}}}
}

func (w *StmtWith) With(name string, builder builder) *StmtWith {
	w.withs = append(w.withs, &with{name: name, builder: builder})

	return w
}

func (w *StmtWith) Select(column ...string) *StmtSelect {
	return newStmtSelect(w.conn.read, w.withs, column)
}

func (w *StmtWith) Insert() *StmtInsert {
	return newStmtInsert(w.conn.write, w.withs)
}

func (w *StmtWith) Update(table string) *StmtUpdate {
	return newStmtUpdate(w.conn.write, w.withs, table)
}

func (w *StmtWith) Delete() *StmtDelete {
	return newStmtDelete(w.conn.write, w.withs)
}

func (w *StmtWith) Build() (string, error) {

	withs, err := w.withs.Build()
	if err != nil {
		return "", err
	}

	// query
	query := fmt.Sprintf("WITH %s", withs)

	return query, nil
}

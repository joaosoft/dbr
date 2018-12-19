package dbr

import (
	"fmt"
)

type StmtWith struct {
	withs       withs
	isRecursive bool

	conn *connections
}

func newStmtWith(conn *connections, name string, isRecursive bool, builder builder) *StmtWith {
	return &StmtWith{conn: conn, withs: withs{&with{name: name, builder: builder}}, isRecursive: isRecursive}
}

func (w *StmtWith) With(name string, builder builder) *StmtWith {
	w.withs = append(w.withs, &with{name: name, builder: builder})

	return w
}

func (w *StmtWith) Select(column ...string) *StmtSelect {
	return newStmtSelect(w.conn.read, w, column)
}

func (w *StmtWith) Insert() *StmtInsert {
	return newStmtInsert(w.conn.write, w)
}

func (w *StmtWith) Update(table string) *StmtUpdate {
	return newStmtUpdate(w.conn.write, w, table)
}

func (w *StmtWith) Delete() *StmtDelete {
	return newStmtDelete(w.conn.write, w)
}

func (w *StmtWith) Build() (string, error) {

	if len(w.withs) == 0 {
		return "", nil
	}

	withs, err := w.withs.Build()
	if err != nil {
		return "", err
	}

	var recursive string
	if w.isRecursive {
		recursive = "RECURSIVE "
	}

	// query
	query := fmt.Sprintf("WITH %s%s", recursive, withs)

	return query, nil
}

package dbr

type condition struct {
	query  string
	values []interface{}
}

type set struct {
	column string
	value  interface{}
}

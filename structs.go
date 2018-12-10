package dbr

type condition struct {
	query  string
	values values
}

type set struct {
	column string
	value  interface{}
}

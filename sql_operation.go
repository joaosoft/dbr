package dbr

type SqlOperation string

const (
	ExecuteOperation SqlOperation = "EXECUTE"
	SelectOperation  SqlOperation = "SELECT"
	InsertOperation  SqlOperation = "INSERT"
	UpdateOperation  SqlOperation = "UPDATE"
	DeleteOperation  SqlOperation = "DELETE"
)

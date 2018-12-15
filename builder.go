package dbr

type builder interface {
	Build() (string, error)
}

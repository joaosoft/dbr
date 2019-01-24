package dbr

type loadOption string

const (
	loadOptionDefault loadOption = "db"
	loadOptionRead    loadOption = "db.read"
	loadOptionWrite   loadOption = "db.write"
)

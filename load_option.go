package dbr

type loadOption string

const (
	loadOptionDefault loadOption = "db"
	loadOptionRead    loadOption = "db.Read"
	loadOptionWrite   loadOption = "db.Write"
)

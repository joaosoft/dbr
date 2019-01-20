package dbr

import (
	"database/sql"
	"time"

	"sync"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Dbr struct {
	conections *connections

	config        *DbrConfig
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
}

type connections struct {
	read     *db
	write    *db
	duration time.Duration
}

type db struct {
	database
	dialect  dialect
	duration time.Duration
}

// New ...
func New(options ...DbrOption) (*Dbr, error) {
	config, simpleConfig, err := NewConfig()

	dbr := &Dbr{
		pm:     manager.NewManager(manager.WithRunInBackground(true)),
		config: &DbrConfig{},
	}

	if err == nil {
		dbr.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Dbr.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	if dbr.isLogExternal {
		dbr.pm.Reconfigure(manager.WithLogger(log))
	}

	dbr.Reconfigure(options...)

	// connect to database
	if dbr.config.Db != nil {
		dbCon := manager.NewSimpleDB(dbr.config.Db)
		if err := dbCon.Start(nil); err != nil {
			return nil, err
		}
		dbr.pm.AddDB("db", dbCon)

		db := &db{database: dbCon.Get(), dialect: newDialect(dbr.config.Db.Driver)}
		dbr.conections = &connections{read: db, write: db}
	} else {
		dbReadCon := manager.NewSimpleDB(dbr.config.ReadDb)
		if err := dbReadCon.Start(nil); err != nil {
			return nil, err
		}
		dbr.pm.AddDB("db-read", dbReadCon)
		dbRead := &db{database: dbReadCon.Get(), dialect: newDialect(dbr.config.ReadDb.Driver)}

		dbWriteCon := manager.NewSimpleDB(dbr.config.WriteDb)
		if err := dbWriteCon.Start(nil); err != nil {
			return nil, err
		}
		dbr.pm.AddDB("db-write", dbWriteCon)
		dbWrite := &db{database: dbReadCon.Get(), dialect: newDialect(dbr.config.WriteDb.Driver)}

		dbr.conections = &connections{read: dbRead, write: dbWrite}
	}

	return dbr, nil
}

func (dbr *Dbr) Select(column ...string) *StmtSelect {
	return newStmtSelect(dbr.conections.read, &StmtWith{}, column)
}

func (dbr *Dbr) Insert() *StmtInsert {
	return newStmtInsert(dbr.conections.write, &StmtWith{})
}

func (dbr *Dbr) Update(table string) *StmtUpdate {
	return newStmtUpdate(dbr.conections.write, &StmtWith{}, table)
}

func (dbr *Dbr) Delete() *StmtDelete {
	return newStmtDelete(dbr.conections.write, &StmtWith{})
}

func (dbr *Dbr) Execute(query string) *StmtExecute {
	return newStmtExecute(dbr.conections.write, query)
}

func (dbr *Dbr) With(name string, builder builder) *StmtWith {
	return newStmtWith(dbr.conections, name, false, builder)
}

func (dbr *Dbr) UseOnlyWrite(name string, builder builder) *StmtWith {
	return newStmtWith(&connections{read: dbr.conections.write, write: dbr.conections.write}, name, false, builder)
}

func (dbr *Dbr) UseOnlyRead(name string, builder builder) *StmtWith {
	return newStmtWith(&connections{read: dbr.conections.read, write: dbr.conections.read}, name, false, builder)
}

func (dbr *Dbr) WithRecursive(name string, builder builder) *StmtWith {
	return newStmtWith(dbr.conections, name, true, builder)
}

func (dbr *Dbr) Begin() (*Transaction, error) {
	tx, err := dbr.conections.write.database.(*sql.DB).Begin()
	if err != nil {
		return nil, err
	}

	return newTransaction(&db{database: tx, dialect: dbr.conections.write.dialect}), nil
}

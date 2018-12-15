package dbr

import (
	"database/sql"
	"fmt"
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
	dbr := &Dbr{
		pm: manager.NewManager(manager.WithRunInBackground(true)),
	}

	if dbr.isLogExternal {
		dbr.pm.Reconfigure(manager.WithLogger(log))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		dbr.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Dbr.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	dbr.config = &appConfig.Dbr

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
	return newStmtSelect(dbr.conections.read, nil, column)
}

func (dbr *Dbr) Insert() *StmtInsert {
	return newStmtInsert(dbr.conections.write, nil)
}

func (dbr *Dbr) Update(table string) *StmtUpdate {
	return newStmtUpdate(dbr.conections.write, nil, table)
}

func (dbr *Dbr) Delete() *StmtDelete {
	return newStmtDelete(dbr.conections.write, nil)
}

func (dbr *Dbr) Execute(query string) *StmtExecute {
	return newStmtExecute(dbr.conections.write, query)
}

func (dbr *Dbr) With(name string, builder builder) *StmtWith {
	return newStmtWith(dbr.conections, name, builder)
}

func (dbr *Dbr) Begin() (*Transaction, error) {
	tx, err := dbr.conections.write.database.(*sql.DB).Begin()
	if err != nil {
		return nil, err
	}

	return newTransaction(&db{database: tx, dialect: dbr.conections.write.dialect}), nil
}

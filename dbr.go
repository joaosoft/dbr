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
	conections *Connections

	config        *DbrConfig
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
}

type Connections struct {
	Read     *Db
	Write    *Db
	Duration time.Duration
}

type Db struct {
	Database
	dialect  Dialect
	Duration time.Duration
}

// NewDbr ...
func NewDbr(options ...DbrOption) (*Dbr, error) {
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

		db := &Db{Database: dbCon.Get(), dialect: newDialect(dbr.config.Db.Driver)}
		dbr.conections = &Connections{Read: db, Write: db}
	} else {
		dbReadCon := manager.NewSimpleDB(dbr.config.ReadDb)
		if err := dbReadCon.Start(nil); err != nil {
			return nil, err
		}
		dbr.pm.AddDB("db-read", dbReadCon)
		dbRead := &Db{Database: dbReadCon.Get(), dialect: newDialect(dbr.config.ReadDb.Driver)}

		dbWriteCon := manager.NewSimpleDB(dbr.config.WriteDb)
		if err := dbWriteCon.Start(nil); err != nil {
			return nil, err
		}
		dbr.pm.AddDB("db-write", dbWriteCon)
		dbWrite := &Db{Database: dbReadCon.Get(), dialect: newDialect(dbr.config.WriteDb.Driver)}

		dbr.conections = &Connections{Read: dbRead, Write: dbWrite}
	}

	return dbr, nil
}

func (dbr *Dbr) Select(column ...string) *StmtSelect {
	return newStmtSelect(dbr.conections.Read, column)
}

func (dbr *Dbr) Insert() *StmtInsert {
	return newStmtInsert(dbr.conections.Write)
}

func (dbr *Dbr) Update(table string) *StmtUpdate {
	return newStmtUpdate(dbr.conections.Write, table)
}

func (dbr *Dbr) Delete() *StmtDelete {
	return newStmtDelete(dbr.conections.Write)
}

func (dbr *Dbr) Begin() (*Transaction, error) {
	tx, err := dbr.conections.Write.Database.(*sql.DB).Begin()
	if err != nil {
		return nil, err
	}

	return newTransaction(&Db{Database: tx, dialect: dbr.conections.Write.dialect}), nil
}

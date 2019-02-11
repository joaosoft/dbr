package dbr

import (
	"database/sql"
	"sync"
	"time"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Dbr struct {
	Connections *connections

	config        *DbrConfig
	logger        logger.ILogger
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
}

type connections struct {
	Read  *db
	Write *db
}

type db struct {
	database
	Dialect dialect
}

func NewDb(database database, dialect dialect) *db {
	return &db{
		database: database,
		Dialect:  dialect,
	}
}

// New ...
func New(options ...DbrOption) (*Dbr, error) {
	config, simpleConfig, err := NewConfig()

	service := &Dbr{
		pm:     manager.NewManager(manager.WithRunInBackground(true)),
		logger: logger.NewLogDefault("dbr", logger.WarnLevel),
		config: config.Dbr,
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Dbr != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Dbr.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	service.Reconfigure(options...)

	// connect to database
	if service.config != nil && service.Connections == nil {
		if service.config.Db != nil {
			dbCon := service.pm.NewSimpleDB(service.config.Db)
			if err := dbCon.Start(nil); err != nil {
				return nil, err
			}
			service.pm.AddDB("db", dbCon)

			db := &db{database: dbCon.Get(), Dialect: NewDialect(service.config.Db.Driver)}
			service.Connections = &connections{Read: db, Write: db}
		} else if service.config.ReadDb != nil && service.config.WriteDb != nil {
			dbReadCon := service.pm.NewSimpleDB(service.config.ReadDb)
			if err := dbReadCon.Start(nil); err != nil {
				return nil, err
			}
			service.pm.AddDB("db-Read", dbReadCon)
			dbRead := &db{database: dbReadCon.Get(), Dialect: NewDialect(service.config.ReadDb.Driver)}

			dbWriteCon := service.pm.NewSimpleDB(service.config.WriteDb)
			if err := dbWriteCon.Start(nil); err != nil {
				return nil, err
			}
			service.pm.AddDB("db-Write", dbWriteCon)
			dbWrite := &db{database: dbReadCon.Get(), Dialect: NewDialect(service.config.WriteDb.Driver)}

			service.Connections = &connections{Read: dbRead, Write: dbWrite}
		}
	}

	return service, nil
}

func (dbr *Dbr) Select(column ...interface{}) *StmtSelect {
	return newStmtSelect(dbr, dbr.Connections.Read, &StmtWith{}, column)
}

func (dbr *Dbr) Insert() *StmtInsert {
	return newStmtInsert(dbr, dbr.Connections.Write, &StmtWith{})
}

func (dbr *Dbr) Update(table string) *StmtUpdate {
	return newStmtUpdate(dbr, dbr.Connections.Write, &StmtWith{}, table)
}

func (dbr *Dbr) Delete() *StmtDelete {
	return newStmtDelete(dbr, dbr.Connections.Write, &StmtWith{})
}

func (dbr *Dbr) Execute(query string) *StmtExecute {
	return newStmtExecute(dbr, dbr.Connections.Write, query)
}

func (dbr *Dbr) With(name string, builder builder) *StmtWith {
	return newStmtWith(dbr, dbr.Connections, name, false, builder)
}

func (dbr *Dbr) UseOnlyWrite(name string, builder builder) *StmtWith {
	return newStmtWith(dbr, &connections{Read: dbr.Connections.Write, Write: dbr.Connections.Write}, name, false, builder)
}

func (dbr *Dbr) UseOnlyRead(name string, builder builder) *StmtWith {
	return newStmtWith(dbr, &connections{Read: dbr.Connections.Read, Write: dbr.Connections.Read}, name, false, builder)
}

func (dbr *Dbr) WithRecursive(name string, builder builder) *StmtWith {
	return newStmtWith(dbr, dbr.Connections, name, true, builder)
}

func (dbr *Dbr) Begin() (*Transaction, error) {
	startTime := time.Now()
	tx, err := dbr.Connections.Write.database.(*sql.DB).Begin()
	if err != nil {
		return nil, err
	}

	return newTransaction(dbr, &db{database: tx, Dialect: dbr.Connections.Write.Dialect}, startTime), nil
}

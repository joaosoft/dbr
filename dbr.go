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
	logger        logger.ILogger
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
	if service.config != nil {
		if service.config.Db != nil {
			dbCon := service.pm.NewSimpleDB(service.config.Db)
			if err := dbCon.Start(nil); err != nil {
				return nil, err
			}
			service.pm.AddDB("Db", dbCon)

			db := &db{database: dbCon.Get(), dialect: newDialect(service.config.Db.Driver)}
			service.conections = &connections{read: db, write: db}
		} else if service.config.ReadDb != nil && service.config.WriteDb != nil {
			dbReadCon := service.pm.NewSimpleDB(service.config.ReadDb)
			if err := dbReadCon.Start(nil); err != nil {
				return nil, err
			}
			service.pm.AddDB("Db-read", dbReadCon)
			dbRead := &db{database: dbReadCon.Get(), dialect: newDialect(service.config.ReadDb.Driver)}

			dbWriteCon := service.pm.NewSimpleDB(service.config.WriteDb)
			if err := dbWriteCon.Start(nil); err != nil {
				return nil, err
			}
			service.pm.AddDB("Db-write", dbWriteCon)
			dbWrite := &db{database: dbReadCon.Get(), dialect: newDialect(service.config.WriteDb.Driver)}

			service.conections = &connections{read: dbRead, write: dbWrite}
		}
	}

	return service, nil
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

package dbr

import (
	"fmt"

	manager "github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Dbr *DbrConfig `json:"dbr"`
}

// DbrConfig ...
type DbrConfig struct {
	Db      *manager.DBConfig `json:"db"`
	ReadDb  *manager.DBConfig `json:"read_db"`
	WriteDb *manager.DBConfig `json:"write_db"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() *DbrConfig {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())
	}

	return appConfig.Dbr
}

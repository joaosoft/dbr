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
	Db      *manager.DBConfig `json:"Db"`
	ReadDb  *manager.DBConfig `json:"read_db"`
	WriteDb *manager.DBConfig `json:"write_db"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}

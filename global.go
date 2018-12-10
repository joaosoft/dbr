package dbr

import logger "github.com/joaosoft/logger"

var log = logger.NewLogDefault("dbr", logger.InfoLevel)
var templates = make(map[string][]byte)

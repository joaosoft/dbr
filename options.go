package dbr

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// DbrOption ...
type DbrOption func(client *Dbr)

// Reconfigure ...
func (dbr *Dbr) Reconfigure(options ...DbrOption) {
	for _, option := range options {
		option(dbr)
	}
}

// WithConfiguration ...
func WithConfiguration(config *DbrConfig) DbrOption {
	return func(dbr *Dbr) {
		dbr.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) DbrOption {
	return func(dbr *Dbr) {
		dbr.logger = logger
		dbr.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) DbrOption {
	return func(dbr *Dbr) {
		dbr.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) DbrOption {
	return func(dbr *Dbr) {
		dbr.pm = mgr
	}
}

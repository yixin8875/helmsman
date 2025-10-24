// Package database provides database client initialization.
package database

import (
	"strings"
	"sync"

	"github.com/go-dev-frame/sponge/pkg/sgorm"

	"helmsman/internal/config"
)

var (
	gdb     *sgorm.DB
	gdbOnce sync.Once

	ErrRecordNotFound = sgorm.ErrRecordNotFound
)

// InitDB connect database
func InitDB() {
	dbDriver := config.Get().Database.Driver
	switch strings.ToLower(dbDriver) {
	case sgorm.DBDriverSqlite:
		gdb = InitSqlite()
	default:
		panic("InitDB error, please modify the correct 'database' configuration at yaml file. " +
			"Refer to https://helmsman/blob/main/configs/helmsman.yml#L85")
	}
}

// GetDB get db
func GetDB() *sgorm.DB {
	if gdb == nil {
		gdbOnce.Do(func() {
			InitDB()
		})
	}

	return gdb
}

// CloseDB close db
func CloseDB() error {
	return sgorm.CloseDB(gdb)
}

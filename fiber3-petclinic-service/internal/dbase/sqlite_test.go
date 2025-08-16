package dbase

import (
	"fiber3-petclinic-service/config/app"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqliteConnectSuccess(t *testing.T) {
	db, _ := SqliteConnect()

	assert.NotNil(t, db)

	sqlDB, _ := db.DB()
	assert.NotNil(t, sqlDB)

	assert.Equal(t, int64(app.Config.Database.MaxIdleConns), sqlDB.Stats().MaxIdleClosed)
	assert.Equal(t, app.Config.Database.MaxOpenConns, sqlDB.Stats().MaxOpenConnections)
}

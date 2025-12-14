package dbase

import (

	"gin-petclinic-service/config/app"
	"gin-petclinic-service/config/app"


	"github.com/stretchr/testify/assert"
)

func TestSqliteConnectSuccess(t *testing.T) {
	db, _ := Connect()

	assert.NotNil(t, db)

	sqlDB, _ := db.DB()
	assert.NotNil(t, sqlDB)

	assert.Equal(t, int64(app.Config.Database.MaxIdleConns), sqlDB.Stats().MaxIdleClosed)
	assert.Equal(t, app.Config.Database.MaxOpenConns, sqlDB.Stats().MaxOpenConnections)
}

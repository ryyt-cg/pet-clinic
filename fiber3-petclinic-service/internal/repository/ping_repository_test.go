package repository

import (
	"fiber3-petclinic-service/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRepository_Ping(t *testing.T) {
	gdb, mock := test.NewMockPostgresDB()
	mock.ExpectQuery("SELECT 1").WithArgs().WillReturnError(nil)
	//if sqlErr != nil {
	//	t.Error(sqlErr)
	//	return
	//}

	pingRepo := NewPingRepository(gdb)
	err := pingRepo.Ping()
	assert.Equal(t, err, nil)
}

package test

import (
	"log"
	"testing"

	"github.com/gorm_use_cases/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	dbcon, err := db.NewDB(sqlite.Open(":memory:"))

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(dbcon); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	dbcon.Exec("PRAGMA foreign_keys = ON")

	return dbcon

}

func TestCreate(t *testing.T) {

	dbcon := GetDB()

	_, err := db.Create(dbcon)
	require.NoError(t, err)

}

func TestRead(t *testing.T) {
	dbcon := GetDB()

	w, err := db.Create(dbcon)
	require.NoError(t, err)

	w2, err := db.Read(dbcon)
	require.NoError(t, err)

	assert.Equal(t, w.WalletID, w2.WalletID)
}

func TestUpdate(t *testing.T) {

	dbcon := GetDB()
	w, err := db.Create(dbcon)
	require.NoError(t, err)

	assert.Equal(t, w.WalletID, "1")
	assert.Equal(t, w.Balance, int64(100))

	w2, err := db.Update(dbcon)
	require.NoError(t, err)

	assert.Equal(t, w.WalletID, w2.WalletID)
}

func TestDelete(t *testing.T) {
	dbcon := GetDB()
	_, err := db.Create(dbcon)
	require.NoError(t, err)

	_, err = db.Delete(dbcon)
	require.NoError(t, err)

	_, err = db.Read(dbcon)

	require.Error(t, err)
	require.Equal(t, err, gorm.ErrRecordNotFound)

}

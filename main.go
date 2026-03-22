package main

import (
	"fmt"
	"log"

	"github.com/gorm_use_cases/db"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	databaseType = "sqlite"
	//database = "mysql"
	dbUser   = "root"
	dbPasswd = "awesome"
	dbName   = "wallet"
)

func main() {

	var dl gorm.Dialector

	switch databaseType {
	//Base de datos sqlite
	case "sqlite":
		dl = sqlite.Open(fmt.Sprintf("%s.db", dbName))

	//Base de datos mySql
	case "mysql":
		dl = mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser,
			dbPasswd,
			dbName,
		),
		)

	//Default
	default:
		log.Fatal("unsupported database")
	}

	//Conecxion a base de datos
	dbcon, err := db.NewDB(dl)

	if err != nil {
		log.Fatal("error: ", err.Error())
	}

	//Migracion automatica
	err = db.AutoMigrate(dbcon)
	if err != nil {
		log.Fatal("error trying to automigrate, err:", err.Error())
	}

	//Uso de operaciones CRUD
	//Create
	//Crea wallet con WalletID = 1
	_, err = db.Create(dbcon)

	//Read
	//Lee wallet con WalletID = 1
	_, err = db.Read(dbcon)
	if err != nil {
		log.Fatal("error trying to read record, err:", err.Error())
	}

	//Update
	//Actualiza wallet con WalletID = 1
	_, err = db.Update(dbcon)
	if err != nil {
		log.Fatal("error trying to update record, err:", err.Error())
	}

	//Delete
	//Borra wallet con WalletID = 1
	_, err = db.Delete(dbcon)
	if err != nil {
		log.Fatal("error trying to delete record, err:", err.Error())
	}

	//Uso de query builder
	//Crea wallet con WalletID =1
	//Busca busca wallet con balance mayor a 500 -> obtiene un registro
	//Busca busca wallet con balance mayor a 5000 -> no obtiene registro
	//Borra wallet con WalletID = 1
	err = db.WhereQueryBuilder(dbcon)
	if err != nil {
		log.Fatal("error trying to use querybuilder", err.Error())
	}

	//Uso de query builder para actualizar registro
	//Crea wallet con WalletID =1
	//Actualiza wallet
	//Borra wallet para evitar conflictos con otros metodos
	err = db.UpdateQueryBuilder(dbcon)
	if err != nil {
		log.Fatal("error trying to use querybuilder", err.Error())
	}

	//Uso de relaciones entre entidades

	//Crea wallet con WalletID = 1
	//Crea transaction con ID = 1 asociado al Wallet con WalletID = 1
	//Crea transaction con ID = 2 asociado al Wallet con WalletID = 1
	//Hace preload de Transacciones
	//Borra registro para evitar conflictos con otras funciones
	err = db.CreateTransaction(dbcon)
	if err != nil {
		log.Fatal("error trying to use create transaction")
	}

	//Raw sql execution
	//Crea wallet con WalletID = 1
	//Lee el wallet con WalletID =1 usando raw SQL
	err = db.RawSQLExecution(dbcon)
	if err != nil {
		log.Fatal("error trying to execute raw sql")
	}

	//SQL transactions
	//Crea wallet con WalletID = 1
	//Crea Transaccion con ID = 1 asociado a wallet con WalletID = 1
	//Crea Transaccion con ID = 2 asociado a wallet con WalletID = 1
	//Si falla guardando alguno de los registros hace rollback
	err = db.SqlTransactions(dbcon)
	if err != nil {
		log.Fatal("error trying to use sql transactions")
	}
}

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
	_, err = db.Create(dbcon)

	//Read
	_, err = db.Read(dbcon)
	if err != nil {
		log.Fatal("error trying to read record, err:", err.Error())
	}

	//Update
	_, err = db.Update(dbcon)
	if err != nil {
		log.Fatal("error trying to update record, err:", err.Error())
	}

	//Delete
	_, err = db.Delete(dbcon)
	if err != nil {
		log.Fatal("error trying to delete record, err:", err.Error())
	}

	//Uso de query builder
	err = db.WhereQueryBuilder(dbcon)
	if err != nil {
		log.Fatal("error trying to use querybuilder", err.Error())
	}

	//Uso de query builder para actualizar registro
	err = db.UpdateQueryBuilder(dbcon)
	if err != nil {
		log.Fatal("error trying to use querybuilder", err.Error())
	}

	//Uso de relaciones entre entidades

	err = db.CreateTransaction(dbcon)
	if err != nil {
		log.Fatal("error trying to use create transaction")
	}

	//Raw sql execution
	err = db.RawSQLExecution(dbcon)
	if err != nil {
		log.Fatal("error trying to execute raw sql")
	}

	//SQL transactions
	err = db.SqlTransactions(dbcon)
	if err != nil {
		log.Fatal("error trying to use sql transactions")
	}
}

package db

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(dl gorm.Dialector) (*gorm.DB, error) {

	db, err := gorm.Open(dl, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA foreign_keys = ON")

	return db, nil

}

func AutoMigrate(db *gorm.DB) error {
	log.Println("runing automigrations")
	return db.AutoMigrate(
		&WalletModel{},
		&TransactionModel{},
	)
}

func Create(db *gorm.DB) (*WalletModel, error) {
	w := WalletModel{
		WalletID: "1",
		Balance:  100,
		Currency: "USD",
	}
	err := db.Create(&w).Error

	if err != nil {
		return nil, err
	}

	return &w, nil

}

func Read(db *gorm.DB) (*WalletModel, error) {
	var w WalletModel
	err := db.First(&w, "1").Error
	if err != nil {
		return nil, err
	}

	log.Println(w)
	return &w, nil

}

func Update(db *gorm.DB) (*WalletModel, error) {
	w := WalletModel{
		WalletID: "1",
		Balance:  1000,
		Currency: "USD",
	}

	err := db.Save(&w).Error
	if err != nil {
		return nil, err
	}
	log.Println(w)
	return &w, nil
}

func Delete(db *gorm.DB) (*WalletModel, error) {
	w := WalletModel{
		WalletID: "1",
		Balance:  0,
		Currency: "USD",
	}
	err := db.Delete(&w).Error
	if err != nil {
		return nil, err
	}

	log.Println(w)
	return &w, nil

}

func WhereQueryBuilder(db *gorm.DB) error {
	//Crea registro
	err := db.Create(&WalletModel{
		WalletID: "1",
		Balance:  1000,
		Currency: "USD",
	}).Error

	if err != nil {
		return err
	}

	var w WalletModel

	//Obtine un registros cuyo balance sea mayor a 500
	err = db.Where("balance > ?", 500).Find(&w).Error
	if err != nil {
		return err
	}

	log.Println(w)

	w = WalletModel{}
	//No optine registro debido a que el balance es menor a 5000
	err = db.Where("balance > ?", 5000).Find(&w).Error
	if err != nil {
		log.Fatal("error trying to read record, err:", err.Error())
	}

	w = WalletModel{
		WalletID: "1",
	}

	//se borra el registro
	err = db.Delete(&w).Error
	if err != nil {
		return err
	}

	log.Println(w)
	return nil
}

func UpdateQueryBuilder(db *gorm.DB) error {

	w := &WalletModel{
		WalletID: "1",
		Balance:  100,
		Currency: "USD",
	}
	err := db.Create(w).Error

	if err != nil {
		return err
	}

	w = &WalletModel{
		WalletID: "1",
		Balance:  1000,
		Currency: "USD",
	}

	err = db.Model(&WalletModel{}).
		Where("wallet_id = ?", w.WalletID).
		Select("*").
		Updates(w).Error

	if err != nil {
		return err
	}

	//Borra wallet para evitar conflictos con otros metodos
	err = db.Delete(&w).Error
	if err != nil {
		return err
	}

	return nil
}

func CreateTransaction(db *gorm.DB) error {

	//Se crea el wallet asociado a la transaccion
	w := WalletModel{
		WalletID: "1",
		Balance:  1000,
		Currency: "USD",
	}
	err := db.Create(&w).Error

	if err != nil {
		return err
	}

	//Se crea la transaccion 1
	t := TransactionModel{
		ID:           "1",
		WalletID:     "1",
		Amount:       500,
		CurrencyCode: "USD",
		Type:         Deposit,
	}

	err = db.Create(&t).Error

	if err != nil {
		return err
	}

	//Se crea la transaccion 2
	t = TransactionModel{
		ID:           "2",
		WalletID:     "1",
		Amount:       500,
		CurrencyCode: "USD",
		Type:         Deposit,
	}

	err = db.Create(&t).Error

	if err != nil {
		return err
	}

	//Uso de preload para obtener las transacciones
	//asociadas

	err = db.Preload("Transactions").Find(&w).Error
	if err != nil {
		return err
	}

	log.Println(w)

	//Borra registro para evitar conflictos con otras funciones
	//Borra transacciones asociadas en cascada
	err = db.Delete(&w).Error
	if err != nil {
		return err
	}

	return nil

}

func RawSQLExecution(db *gorm.DB) error {
	//Se crea el wallet
	w := WalletModel{
		WalletID: "1",
		Balance:  1000,
		Currency: "USD",
	}
	err := db.Create(&w).Error

	if err != nil {
		return err
	}

	w = WalletModel{}
	err = db.Raw("SELECT * FROM wallet_models WHERE wallet_id = ?", "1").Scan(&w).Error
	if err != nil {
		return err
	}
	log.Println(w)

	err = db.Delete(&w).Error
	if err != nil {
		return err
	}

	return nil
}

func SqlTransactions(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//Se crea el wallet asociado a la transaccion
		w := WalletModel{
			WalletID: "1",
			Balance:  1000,
			Currency: "USD",
		}
		err := tx.Create(&w).Error

		if err != nil {
			return err //rollback
		}

		//Se crea la transaccion 1
		t := TransactionModel{
			ID:           "1",
			WalletID:     "1",
			Amount:       500,
			CurrencyCode: "USD",
			Type:         Deposit,
		}

		err = tx.Create(&t).Error

		if err != nil {
			return err //rollback
		}

		//Se crea la transaccion 2
		t = TransactionModel{
			ID:           "2",
			WalletID:     "1",
			Amount:       500,
			CurrencyCode: "USD",
			Type:         Deposit,
		}

		err = tx.Create(&t).Error

		if err != nil {
			return err //rollback
		}

		//Retorno de error comentado, habilite para observar el funcionamiento de rollback
		//return errors.New("rollback error")

		return nil
	})

	if err != nil {
		return err
	}

	err = db.Delete(&WalletModel{WalletID: "1"}).Error
	if err != nil {
		return err
	}
	return nil
}

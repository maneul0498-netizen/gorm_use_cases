# 📌 Resumen principales caracteristicas de GORM usadas

Demostracion simple de uso de gorm con MySQL y SQLite
---

## 🛠️ Tecnologías

- Lenguaje: Go, GORM
- Base de datos: MySQL / SQLite

---

Para correr el proyecto haga uso de la variable databaseType asignandole el valor "sqlite". Si desea
correr el proyecto con base de datos mysql haga uso de las variables dbUser, dbPasswd y dbName y asigne
el valor "mysql" a la variable databaseType
```go
var (
	databaseType = "sqlite"
	//database = "mysql"
	dbUser   = "root"
	dbPasswd = "awesome"
	dbName   = "wallet"
)
```

---

#### Se agregaron pruebas de integracion con sqlite a algunas de las funciones implementadas para correrlos ejecute el sigueinte comando en la raiz del proyecto. 
```bash
go test -v ./...
```
#### Nota: para cada funcion de test se usa una instancia de conexion diferente, (podria usarse un singleton pero no se usa a fin de que cada una de las pruebas se ejecute de forma aislada)
---

# 📌 Resumen principales caracteristicas de GORM usadas

## Uso de orm GORM con bases de datos SQLite y MySql

#### Uso de migraciones automaticas:
```go
type WalletModel struct {
	WalletID  string `gorm:"primaryKey"`
	Balance   int64  `gorm:"not null;default:0"` // en centavos (nunca float)
	Currency  string `gorm:"type:varchar(10);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Transactions []TransactionModel `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
db.AutoMigrate(&WalletModel{})
```

---

#### Uso de operaciones CRUD:
```go
db.Create(&wallet)    // Crear
db.First(&wallet, 1)  // Leer
db.Save(&wallet)      // Actualizar
db.Delete(&wallet)    // Eliminar
```

---

#### Uso de Query Builder:
```go
db.Where("balance > ?", 500).Find(&wallet)
```

---

#### Uso de relaciones entre modelos tales como:
- Uno a uno
- Uno a muchos
- Muchos a uno

#### Ejemplo de relacion uno a muchos:
```go
type WalletModel struct {
	WalletID  string `gorm:"primaryKey"`
	Balance   int64  `gorm:"not null;default:0"` // en centavos (nunca float)
	Currency  string `gorm:"type:varchar(10);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Transactions []TransactionModel `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type TransactionModel struct {
	ID           string          `gorm:"primaryKey"`
	WalletID     string          `gorm:"not null;index"`
	Amount       int64           `gorm:"not null"` // en centavos
	CurrencyCode string          `gorm:"not null"`
	Type         TransactionType `gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
```
#### Uso de preload para cargar relaciones:
```go
db.Preload("Transactions").Find(&wallet)
```
---

#### Uso de transacciones. Ver funcion SqlTransactions, contiene un ejemplo mas realista
```go
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err // rollback
    }
    return nil // commit
})
```

---

#### Ejecucion de raw SQL:
```go
db.Raw("SELECT * FROM wallet_models").Scan(&wallet)
```

---

#### Uso de hooks: Los metdos son invocados por gorm al hacer el llamado a Create
```go
func (w *WalletModel) BeforeCreate(tx *gorm.DB) error {
	log.Println("Executing BeforeCreate (walletModel)")
	return nil

}

func (w *WalletModel) AfterCreate(tx *gorm.DB) error {
	log.Println("Executing AfterCreate (walletModel)")
	return nil

}
db.Create(&wallet)
```
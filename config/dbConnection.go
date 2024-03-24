package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/patnaikankit/POS-API/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	godotenv.Load()
	dbHost := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DBNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbUser := os.Getenv("MYSQL_USER")

	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", dbUser, dbPassword, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
	fmt.Println(dbHost)

	if err != nil {
		panic("Db connection failed")
	}

	DB = db
	fmt.Println("Connection successful")

	AutoMigrate(db)
}

func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&models.Cashier{},
		&models.Category{},
		&models.Discount{},
		&models.Order{},
		&models.Payment{},
		&models.Product{},
		&models.PaymentType{},
	)
}

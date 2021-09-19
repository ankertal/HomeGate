package server

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDatabase returns database connection
func GetDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("HOMEGATE_PSQL_USER"), os.Getenv("HOMEGATE_PSQL_PASSWORD"), os.Getenv("HOMEGATE_PSQL_DBNAME"), os.Getenv("HOMEGATE_PSQL_DBPORT"))

	log.Print("Connecting to Postgres DB...")
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Invalid database url")
	}
	sqldb, _ := connection.DB()

	err = sqldb.Ping()
	if err != nil {
		log.Fatal("Database connected")
	}
	fmt.Println("Database connection successful.")
	return connection
}

//create user table in userdb
func InitialMigration() {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.AutoMigrate(User{}, Gate{})
}

// CloseDatabase closes database connection
func CloseDatabase(connection *gorm.DB) {
	sqldb, _ := connection.DB()
	sqldb.Close()
}

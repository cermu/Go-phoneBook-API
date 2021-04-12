package models

import (
	"fmt"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
)

var DBConnection *gorm.DB

// InitDB public function used to create a pointer to database connection
func InitDB() {
	log.Printf("INFO | Initializing database ...")
	var err error

	dbName := utl.ReadConfigs().GetString("DB.NAME")
	dbUser := utl.ReadConfigs().GetString("DB.USER")
	dbPass := utl.ReadConfigs().GetString("DB.PASS")
	dbHost := utl.ReadConfigs().GetString("DB.HOST")
	dbPort := utl.ReadConfigs().GetInt("DB.PORT")
	// dbType := utl.ReadConfigs().GetString("DB.TYPE")

	dns := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName, dbPass)
	DBConnection, err = gorm.Open("postgres", dns)

	DBConnection.SingularTable(true)

	if err != nil {
		log.Fatalf("WARNING | Database connection failed with message: %v\n", err.Error())
	}

	log.Printf("INFO | Initializing database \t [OK]")
}

// MigrateDB public function used to perform database migrations.
// Our models will be translated to database tables
func MigrateDB () {
	log.Println("INFO | Running database migrations ...")
	DBConnection.Debug().AutoMigrate(Account{}, Contact{})
	// DBConnection.Debug().AUtoMigrate(...)

	// migrating foreign keys
	DBConnection.Model(&Contact{}).AddForeignKey("account_id", "account(id)", "CASCADE", "CASCADE")
	log.Println("INFO | Database migrations completed")
}

package dbhelper

import (
	"fmt"
	"log"
	"os"
	"time"
	"wallet-service-gin/utils"

	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupDB function to setup the db and check the db connection
func SetupDB() (*gorm.DB, error) {
	// Config to show log of DB Query.
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Include params in the SQL log
			Colorful:                  true,        // Enable color
		},
	)

	// Open the DB connection
	var err error
	db, err := gorm.Open(mysql.Open(DSN()), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Check the db connection
	pingErr := sqlDB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return nil, pingErr
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(15 * time.Minute)

	log.Printf("=======================")
	log.Printf("MaxIdleConns : %d", 10)
	log.Printf("MaxOpenConns : %d", 100)
	log.Printf("-----------------------")
	log.Printf("DB Connected!")
	log.Printf("=======================")

	return db, nil
}

// DSN make string format of DB config
func DSN() string {
	// DB Config setting for local
	config := &gomysql.Config{
		User:      utils.GetEnvDefault("DBUSER", "root"),
		Passwd:    utils.GetEnvDefault("DBPASS", ""),
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%s", utils.GetEnvDefault("DBHOST", "localhost"), utils.GetEnvDefault("DBPORT", "3306")),
		DBName:    "walletapp",
		ParseTime: true,
	}

	return config.FormatDSN()
}

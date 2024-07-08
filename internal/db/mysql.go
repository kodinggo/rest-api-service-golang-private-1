package db

import (
	"database/sql"
	"fmt"
	"log"
)

func InitMySQLConn() *sql.DB {
	// db environment
	var (
		dbUsername = "root"
		dbPassword = ""
		dbName     = "kodinggo"
		dbHost     = "localhost:3306"
	)

	// prepare connection string
	// charset=utf8mb4 agar dapat menyimpan semua karakter unicode
	// parseTime=true agar dapat diparsing dari timestamp ke tipe time.Time
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true",
		dbUsername,
		dbPassword,
		dbHost,
		dbName)
	connDB, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Panicf("failed to connect server db, error: %s", err.Error())
	}

	return connDB
}

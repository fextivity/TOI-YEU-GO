package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func New(connection *string) (*DB, error) {
	if connection == nil {
		connection = loadDBConfig()
	}

	log.Print(*connection)

	db_handle, err := sql.Open("mysql", *connection)
	if err != nil {
		return nil, err
	}

	db := DB{
		db_handle,
	}

	return &db, nil
}

func loadDBConfig() *string {
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PWD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}

	dsn := cfg.FormatDSN()

	return &dsn
}

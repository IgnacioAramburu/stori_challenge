package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDB() (*sql.DB, error) {
	dsn := "root:satorugojo2208@tcp(localhost:3306)/stori_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}

package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func Connect() error {
	connString := os.Getenv("POSTGRES_CONN_STRING")

	var err error
	Conn, err = sql.Open("postgres", connString)
	if err != nil {
		log.Println("error while opening database connection: ", err)
		return err
	}

	err = Conn.Ping()
	if err != nil {
		log.Println("error while testing database connection: ", err)
		return err
	}

	return nil
}

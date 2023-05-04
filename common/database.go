package common

import (
	"log"
	"database/sql"
	_ "github.com/alexbrainman/odbc"
)

var DB *sql.DB

func InitDB() *sql.DB {
	connStr := "DRIVER={Microsoft Access Driver (*.mdb, *.accdb)};FIL={MS Access};DBQ=D:\\AccessSql\\User.accdb"
	db, err := sql.Open("odbc", connStr)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	return db
}

func GetDB() *sql.DB{
	return DB
}
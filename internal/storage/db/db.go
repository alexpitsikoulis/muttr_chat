package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	return nil
}

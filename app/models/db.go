package models

import (
	"fmt"
	"log"
	"suger-clickup/pkg/settings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS account (
	user_id serial PRIMARY KEY,
	username VARCHAR (50) UNIQUE NOT NULL,
	password VARCHAR (50) NOT NULL,
	email VARCHAR (255) UNIQUE NOT NULL,
	created_at TIMESTAMP NOT NULL
);


`

var db *sqlx.DB

func Setup() {
	setupDB()
}

func setupDB() {
	conf := settings.Get()
	var err error
	uri := conf.DBConfig.Format()
	fmt.Println("Uri", uri)
	db, err = sqlx.Connect("postgres", conf.DBConfig.Format())
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

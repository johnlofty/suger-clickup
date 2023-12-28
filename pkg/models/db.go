package models

import (
	"fmt"
	"log"
	"suger-clickup/pkg/settings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS orgs (
	org_id serial PRIMARY KEY,
	org_name VARCHAR(255) UNIQUE NOT NULL
);



CREATE TABLE IF NOT EXISTS accounts (
	user_id serial PRIMARY KEY,
	email VARCHAR (255) UNIQUE NOT NULL,
	password VARCHAR (50) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	org_id INT NOT NULL,
	FOREIGN KEY (org_id) REFERENCES orgs(org_id)
);

CREATE TABLE IF NOT EXISTS tickets (
	ticket_id VARCHAR (36) PRIMARY KEY,
	user_id INT NOT NULL,
	org_id INT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	FOREIGN KEY (org_id) REFERENCES orgs(org_id),
	FOREIGN KEY (user_id) REFERENCES accounts(user_id)
);

CREATE TABLE IF NOT EXISTS ticket_assignees (
	ticket_id VARCHAR (36) REFERENCES tickets(ticket_id),
	assignee_id INT REFERENCES accounts(user_id),
	PRIMARY KEY (ticket_id, assignee_id)
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
	db.MustExec(schema)
}

func GetDB() *sqlx.DB {
	return db
}

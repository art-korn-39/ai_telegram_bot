package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Operation struct {
	date     time.Time
	user_id  int64
	username string
	model    string
	request  string
}

func SQL_Connect() {

	// Capture connection properties.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Cfg.DB_host, Cfg.DB_port, Cfg.DB_user, Cfg.DB_password, Cfg.DB_name)

	// Get a database handle.
	var err error
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println("Unsuccessful connection to PostgreSQL!")
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println("Unsuccessful connection to PostgreSQL!")
		log.Fatal(pingErr)
	}

	log.Println("Successful connection to DB " + Cfg.DB_name)

}

func NewSQLOperation(user *UserInfo, request string) Operation {

	return Operation{
		date:     time.Now().UTC().Add(3 * time.Hour),
		user_id:  user.ChatID,
		username: subString(user.Username, 0, 40),
		model:    user.Model,
		request:  request,
	}
}

func SQL_AddOperation(o Operation) {

	if db == nil {
		Logs <- Log{"sql", "lost connection to DB", true}
		return
	}

	Statement := `
	INSERT INTO Operations (date, user_id, username, model, request)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(Statement,
		o.date,
		o.user_id,
		o.username,
		o.model,
		o.request)

	if err != nil {
		Logs <- Log{"sql", err.Error(), true}
		return
	}

}

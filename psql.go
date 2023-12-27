package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type Operation struct {
	data     time.Time
	user_id  int64
	username string
	model    string
	request  string
}

func SQL_Connect() {

	return

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

func NewSQLOperation(user *UserInfo, upd tgbotapi.Update, request string) Operation {

	return Operation{
		data:     time.Now(),
		user_id:  upd.Message.Chat.ID,
		username: upd.Message.From.UserName,
		model:    user.Model,
		request:  request,
	}
}

func SQL_AddOperation(o Operation) {

	return

	if db == nil {
		Logs <- Log{"sql", "lost connection to DB", true}
		return
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	Statement := `
	INSERT INTO Operations (data, user_id, username, model, request)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.Exec(Statement,
		o.data,
		o.user_id,
		o.username,
		o.model,
		o.request)

	tx.Commit()

	if err != nil {
		Logs <- Log{"sql", err.Error(), true}
	}
}

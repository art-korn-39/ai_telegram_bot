package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type Operation struct {
	date     time.Time
	chat_id  int64
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

func SQL_NewOperation(user *UserInfo, request string) Operation {

	return Operation{
		date:     time.Now().UTC().Add(3 * time.Hour),
		chat_id:  user.ChatID,
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
	INSERT INTO Operations (date, chat_id, username, model, request)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(Statement,
		o.date,
		o.chat_id,
		o.username,
		o.model,
		o.request)

	if err != nil {
		Logs <- Log{"sql", err.Error(), true}
		return
	}

}

func SQL_LoadUserStates() {

	if db == nil {
		log.Fatal("sql", "lost connection to DB")
	}

	stmt := `
	select
		user_name, chat_id, model, last_command, input_text, stage 
	from 
		user_states
	`
	rows, err := db.Query(stmt)
	if err != nil {
		Logs <- Log{"sql_load", err.Error(), true}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u UserInfo
		if err := rows.Scan(&u.Username, &u.ChatID, &u.Model,
			&u.LastCommand, &u.InputText, &u.Stage); err != nil {
			Logs <- Log{"sql_load", err.Error(), true}
		}
		ListOfUsers[u.ChatID] = &u
	}
	if err = rows.Err(); err != nil {
		Logs <- Log{"sql_load", err.Error(), true}
		return
	}

	log.Println("Loading user_states complete")

}

func SQL_SaveUserStates() {

	if db == nil {
		log.Printf("[%s] %s", "sql", "lost connection to DB")
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	stmt := `delete from user_states`
	_, err := tx.Exec(stmt)
	if err != nil {
		log.Printf("[%s] %s", "sql", err.Error())
		return
	}

	stmt = `insert into user_states (user_name, chat_id, model, last_command, input_text, stage)
	values ($1, $2, $3, $4, $5, $6)`

	for _, v := range ListOfUsers {
		_, err = tx.Exec(stmt, v.Username, v.ChatID, v.Model, v.LastCommand, v.InputText, v.Stage)
		if err != nil {
			log.Printf("[%s] %s", "sql", err.Error())
			return
		}
	}

	log.Printf("[%s] %s", "sql", "Saving user states done")
	tx.Commit()

}

func SQL_GetInfoOnDate(timestamp time.Time) (result map[string]int, errStr string) {

	if db == nil {
		Logs <- Log{"sql", "lost connection to DB", true}
		return result, "Отсутствует подключение к БД"
	}

	result = map[string]int{}
	var count int

	Statement := `
	select count(distinct username) from operations where date > '$1';
	select count(*), model from operations where date > '$1' group by model;
	`
	Statement = strings.ReplaceAll(Statement, "$1", timestamp.Format(time.DateTime))

	rows, err := db.Query(Statement)
	if err != nil {
		Logs <- Log{"info", err.Error(), true}
		return result, "Ошибка при выполнении запроса к БД"
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return result, err.Error()
		}
	}
	if err = rows.Err(); err != nil {
		return result, err.Error()
	}

	result["users"] = count

	rows.NextResultSet()

	for rows.Next() {
		var model string
		err := rows.Scan(&count, &model)
		switch err {
		case nil:
			result[model] = count
		default:
			return result, err.Error()
		}
	}
	if err = rows.Err(); err != nil {
		return result, err.Error()
	}

	return

}

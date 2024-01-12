package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Operation struct {
	date     time.Time
	chat_id  int64
	username string
	model    string
	class    string
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
		log.Println(err.Error())
		return
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println("Unsuccessful connection to PostgreSQL!")
		log.Println(pingErr.Error())
		return
	}

	log.Println("Successful connection to DB " + Cfg.DB_name)

}

func SQL_NewOperation(user *UserInfo, model, class, request string) Operation {

	return Operation{
		date:     time.Now().UTC().Add(3 * time.Hour),
		chat_id:  user.ChatID,
		username: SubString(user.Username, 0, 40),
		model:    model,
		class:    class,
		request:  request,
	}
}

func SQL_AddOperation(o Operation) {

	if db == nil {
		Logs <- NewLog(nil, "SQL{AddOperation}", Error, "lost connection to DB")
		return
	}

	Statement := `
	INSERT INTO operations (date, chat_id, username, model, class, request)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(Statement,
		o.date,
		o.chat_id,
		o.username,
		o.model,
		o.class,
		o.request)

	if err != nil {
		Logs <- NewLog(nil, "SQL{AddOperation}", Error, err.Error())
		return
	}

}

func SQL_LoadUserStates() {

	if db == nil {
		Logs <- NewLog(nil, "SQL{LoadUserStates}", Error, "lost connection to DB")
		return
	}

	stmt := `
	SELECT
		user_name, chat_id, path, options, tokens_used_gpt, language, requests_today_gen 
	FROM 
		user_states
	`
	rows, err := db.Query(stmt)
	if err != nil {
		Logs <- NewLog(nil, "SQL{LoadUserStates}", Error, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u UserInfo
		var options string
		if err := rows.Scan(
			&u.Username, &u.ChatID,
			&u.Path, &options,
			&u.Tokens_used_gpt, &u.Language,
			&u.Requests_today_gen); err != nil {
			Logs <- NewLog(nil, "SQL{LoadUserStates}", Error, err.Error())
		}
		u.Options = JSONtoMap(options)
		ListOfUsers[u.ChatID] = &u
	}
	if err = rows.Err(); err != nil {
		Logs <- NewLog(nil, "SQL{LoadUserStates}", Error, err.Error())
		return
	}

	Logs <- NewLog(nil, "SQL{LoadUserStates}", Info, "Loading user_states complete")

}

func SQL_SaveUserStates() {

	if db == nil {
		Logs <- NewLog(nil, "SQL{SaveUserStates}", Error, "lost connection to DB")
		return
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	stmt := `delete from user_states`
	_, err := tx.Exec(stmt)
	if err != nil {
		Logs <- NewLog(nil, "SQL{SaveUserStates}", Error, err.Error())
		return
	}

	stmt = `INSERT INTO user_states (user_name, chat_id, path, options, tokens_used_gpt, language, requests_today_gen)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, v := range ListOfUsers {
		optionsJSON := MapToJSON(v.Options)
		_, err = tx.Exec(stmt, v.Username, v.ChatID, v.Path, optionsJSON, v.Tokens_used_gpt, v.Language, v.Requests_today_gen)
		if err != nil {
			Logs <- NewLog(nil, "SQL{SaveUserStates}", Error, err.Error())
			return
		}
	}

	Logs <- NewLog(nil, "SQL{SaveUserStates}", Info, "Saving user states done")

	tx.Commit()

}

func SQL_GetInfoOnDate(timestamp time.Time) (result map[string]int, errStr string) {

	if db == nil {
		Logs <- NewLog(nil, "SQL{Info}", Error, "lost connection to DB")
		return result, "Отсутствует подключение к БД"
	}

	result = map[string]int{}
	var count int

	Statement := `
	SELECT count(distinct username) FROM operations WHERE date > '$1';
	SELECT count(*), model FROM operations WHERE date > '$1' GROUP BY model;
	`
	Statement = strings.ReplaceAll(Statement, "$1", timestamp.Format(time.DateTime))

	rows, err := db.Query(Statement)
	if err != nil {
		Logs <- NewLog(nil, "SQL{Info}", Error, err.Error())
		return nil, err.Error()
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return nil, err.Error()
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err.Error()
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
		return nil, err.Error()
	}

	return

}

func SQL_AddLog(l Log) {

	if db == nil {
		Logs <- NewLog(nil, "SQL{AddLog}", Error, "lost connection to DB")
		return
	}

	stat := `
	INSERT INTO logs (date, chat_id, author, path, level, text)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(stat,
		l.Date,
		l.ChatID,
		l.Author,
		l.Path,
		l.Level,
		l.Text)

	if err != nil {
		Logs <- NewLog(nil, "SQL{AddLog}", Error, err.Error())
		return
	}

}

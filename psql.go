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
		username: user.Username,
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

func startFillSQL() {

	tx, _ := db.Begin()
	defer tx.Rollback()

	Statement := `
	ALTER TABLE operations 
	RENAME COLUMN data TO date;`
	_, err := db.Exec(Statement)
	if err != nil {
		tx.Rollback()
		fmt.Println("tx fail 1")
		return
	}

	Statement = `
	update operations
	set date = date + interval '3 hour';`
	_, err = db.Exec(Statement)
	if err != nil {
		tx.Rollback()
		fmt.Println("tx fail 2")
		return
	}

	//27.12_21.00 G:53 C:14 K:41
	t := time.Date(2023, 12, 27, 21, 0, 0, 0, time.Local)
	models := map[string]int{"gemini": 53, "chatgpt": 14, "kandinsky": 41}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	//27.12_23.59 G:25 C:41 K:60
	t = time.Date(2023, 12, 27, 23, 59, 59, 0, time.Local)
	models = map[string]int{"gemini": 25, "chatgpt": 41, "kandinsky": 60}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	//28.12_04.00 G:20 C:36 K:51
	t = time.Date(2023, 12, 28, 4, 0, 0, 0, time.Local)
	models = map[string]int{"gemini": 20, "chatgpt": 36, "kandinsky": 51}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	//28.12_07.00 G:4 C:0 K:52
	t = time.Date(2023, 12, 28, 7, 0, 0, 0, time.Local)
	models = map[string]int{"gemini": 4, "chatgpt": 0, "kandinsky": 52}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	// 28.12_09.00 G:6 C:16 K:4
	t = time.Date(2023, 12, 28, 9, 0, 0, 0, time.Local)
	models = map[string]int{"gemini": 6, "chatgpt": 16, "kandinsky": 4}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	// 28.12_18.00 G:24 C:59 K:80
	t = time.Date(2023, 12, 28, 18, 0, 0, 0, time.Local)
	models = map[string]int{"gemini": 24, "chatgpt": 59, "kandinsky": 80}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	// 28.12_22.00 G:14 C:34 K:25
	t = time.Date(2023, 12, 28, 22, 0, 0, 0, time.Local)
	models = map[string]int{"gemini": 14, "chatgpt": 34, "kandinsky": 25}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	// 28.12_23.59 G:20 C:4 K:3
	t = time.Date(2023, 12, 28, 23, 59, 59, 0, time.Local)
	models = map[string]int{"gemini": 20, "chatgpt": 4, "kandinsky": 3}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	// 29.12_03.00 G:18 C:6 K:3
	t = time.Date(2023, 12, 29, 3, 0, 0, 0, time.Local)
	models = map[string]int{"gemini": 18, "chatgpt": 6, "kandinsky": 3}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	// 29.12_16.30 G:19 C:52 K:45
	t = time.Date(2023, 12, 29, 16, 30, 0, 0, time.Local)
	models = map[string]int{"gemini": 19, "chatgpt": 52, "kandinsky": 45}
	for model, count := range models {
		for i := 1; i <= count; i++ {
			SQL_AddOperation(Operation{date: t, model: model})
		}
	}

	tx.Commit()

	fmt.Println("DONE")

}

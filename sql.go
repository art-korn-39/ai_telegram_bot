package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// https://habr.com/ru/companies/oleg-bunin/articles/461935/
// https://github.com/jmoiron/sqlx
// https://www.sobyte.net/post/2021-06/sqlx-library-usage-guide/

const sql_LostConnection = "lost connection to DB"
const sql_AddLog = "SQL{AddLog}"

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
	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Println("Unsuccessful connection to PostgreSQL!")
		log.Println(err.Error())
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
		Logs <- NewLog(nil, "SQL{AddOperation}", Error, sql_LostConnection)
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
		Logs <- NewLog(nil, "SQL{LoadUserStates}", Error, sql_LostConnection)
		return
	}

	stmt := `
	SELECT
		user_name, chat_id, path, options, tokens_used_gpt, 
		language, requests_today_gen, requests_today_sdxl, level 
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
			&u.Requests_today_gen, &u.Requests_today_sdxl,
			&u.Level); err != nil {
			Logs <- NewLog(nil, "SQL{LoadUserStates}", Error, err.Error())
		}
		u.Options = JSONtoMap(options)
		ListOfUsers[u.ChatID] = &u
	}
	if err = rows.Err(); err != nil {
		Logs <- NewLog(nil, "SQL{LoadUserStates}", Error, err.Error())
		return
	}

	Logs <- NewLog(nil, "SQL", Info, "Loading user_states complete")

}

func SQL_SaveUserStates() {

	if db == nil {
		Logs <- NewLog(nil, "SQL{SaveUserStates}", Error, sql_LostConnection)
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

	stmt = `INSERT INTO user_states (user_name, chat_id, path, options, 
									 tokens_used_gpt, language, requests_today_gen, requests_today_sdxl, 
									 level)
	VALUES ($1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9)`

	for _, v := range ListOfUsers {
		optionsJSON := MapToJSON(v.Options)
		_, err = tx.Exec(stmt,
			v.Username, v.ChatID, v.Path, optionsJSON,
			v.Tokens_used_gpt, v.Language, v.Requests_today_gen, v.Requests_today_sdxl,
			v.Level)
		if err != nil {
			Logs <- NewLog(nil, "SQL{SaveUserStates}", Error, err.Error())
			return
		}
	}

	Logs <- NewLog(nil, "SQL", Info, "Saving user states done")

	tx.Commit()

}

func SQL_GetInfoOnDate(timestamp time.Time) (result map[string]int, errStr string) {

	if db == nil {
		Logs <- NewLog(nil, "SQL{Info}", Error, sql_LostConnection)
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

func SQL_GetNewUsersForDay(timestamp time.Time) (cnt int, errStr string) {

	if db == nil {
		Logs <- NewLog(nil, "SQL{SQL_GetNewUsersForToday}", Error, sql_LostConnection)
		return 0, "Отсутствует подключение к БД"
	}

	Statement := `
	with first_days_by_users as (select 
		min(date_trunc('day', date)) as date, 
		chat_id 
		from operations 
		group by chat_id)

	select 
	count(*) as cnt  
	from first_days_by_users
	where date = date_trunc('day', $1::timestamp)
	group by date;`

	err := db.Get(&cnt, Statement, timestamp)
	if err != nil {

		if err == sql.ErrNoRows {
			return 0, ""
		} else {
			Logs <- NewLog(nil, "SQL{SQL_GetNewUsersForToday}", Error, err.Error())
			return 0, err.Error()
		}

	}

	return cnt, ""

}

func SQL_AddLog(l Log) {

	if db == nil {
		Logs <- NewLog(nil, sql_AddLog, Error, sql_LostConnection)
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
		Logs <- NewLog(nil, sql_AddLog, Error, err.Error())
		return
	}

}

func SQL_CountOfUserOperations(u *UserInfo) (count int, isErr bool) {

	if db == nil {
		Logs <- NewLog(nil, "SQL{CountOfUserOperations}", Error, sql_LostConnection)
		return 0, true
	}

	stmt := `
	select count(*) as cnt
	from operations 
	where chat_id = $1
	`
	rows, err := db.Query(stmt, u.ChatID)
	if err != nil {
		Logs <- NewLog(nil, "SQL{CountOfUserOperations}", Error, err.Error())
		return 0, true
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			Logs <- NewLog(nil, "SQL{CountOfUserOperations}", Error, err.Error())
			return 0, true
		}
	}
	if err = rows.Err(); err != nil {
		Logs <- NewLog(nil, "SQL{CountOfUserOperations}", Error, err.Error())
		return 0, true
	}

	return count, false

}

// Получает все замеченные chat_id для массовой рассылки в боте
func SQL_GetAllUsers() (users []*UserInfo, isErr bool) {

	if db == nil {
		Logs <- NewLog(nil, "SQL{GetAllUsers}", Error, sql_LostConnection)
		return nil, true
	}

	stmt := `
	WITH chats as 
		(select chat_id 
		from user_states
		union
		select chat_id 
		from logs
		union
		select chat_id 
		from operations)

	select 
		c.chat_id as chat_id, 
		COALESCE(us.language, '') as language
	from chats c
		left join user_states us 
		on c.chat_id = us.chat_id
	order by c.chat_id;`

	err := db.Select(&users, stmt)
	if err != nil {
		Logs <- NewLog(nil, "SQL{GetAllUsers}", Error, err.Error())
		return nil, true
	}

	return users, false

}

// Возвращает длину последней серии дней
// Вернёт 1 если записи были сегодня, а вчера нет
// Если вчера и сегодня не было записей в логах, то chat_id не будет в ключах мапы-результата
func SQL_UserDayStreak(user *UserInfo) (daysByUsers map[int64]int, isErr bool) {

	daysByUsers = map[int64]int{}

	if db == nil {
		Logs <- NewLog(nil, "SQL{UserDayStreak}", Error, sql_LostConnection)
		return daysByUsers, true
	}

	// Если user не указан, то берем все ID из мапы ListOfUsers
	var sliceChatID []int64
	if user == nil {
		sliceChatID = make([]int64, 0, len(ListOfUsers))
		for _, v := range ListOfUsers {
			sliceChatID = append(sliceChatID, v.ChatID)
		}
	} else {
		sliceChatID = make([]int64, 0, 1)
		sliceChatID = append(sliceChatID, user.ChatID)
	}

	Statement := `
		-- все дни использования и предыдущий день
		with tmp1 as (select distinct
		    chat_id::bigint as chat_id,
			date_trunc('day', date) as day,
			date_trunc('day', date) - interval '1 day' as prevDay
			from logs
		where chat_id::bigint = ANY($1)
			order by day),
				
		-- все дни использования и метка начала каждой серии
		tmp2 as (select 
		    t1.chat_id as chat_id,
			t1.day as day, 
			t2.day is null as firstDay
			from tmp1 t1
			left join tmp1 t2 
				on t1.prevDay = t2.day and t1.chat_id = t2.chat_id),
				
		-- первый день самой последней серии использования
		tmp3 as (select 
			max(day) as day,
			chat_id as chat_id
			from (select 
				day, chat_id 
				from tmp2 
				where firstDay = true) as t 
			group by 
			    chat_id),
		
		-- последний день использования бота
		lastDays as ( select max(day) as day, chat_id from tmp1 group by chat_id),
		
		-- были операции за день до момента среза
		prevDayHasLogs as ( 
			select 
			chat_id as chat_id, 
		    day >= ( date_trunc('day', $2::timestamp) - interval '1 day' ) as stat 
			from lastDays )
			
		-- итоговая выборка
		select 
		t.chat_id as chat_id,
		count(*) as days
		from tmp1 t 
		inner join tmp3 filter 
			on t.day >= filter.day and t.chat_id = filter.chat_id
		inner join prevDayHasLogs
			on t.chat_id = prevDayHasLogs.chat_id
		where 
		    prevDayHasLogs.stat = true
		group by 
		    t.chat_id;`

	type userDays struct {
		Days    int64 `db:"days"`
		Chat_id int64 `db:"chat_id"`
	}
	var data []userDays

	err := db.Select(&data, Statement, pq.Array(sliceChatID), MskTimeNow())
	if err != nil {
		fmt.Println(err.Error())
		Logs <- NewLog(nil, "SQL{UserDayStreak}", Error, err.Error())
		return daysByUsers, true
	}

	for _, v := range data {
		daysByUsers[v.Chat_id] = int(v.Days)
	}

	return daysByUsers, false

}

func SQL_GetFirstDate(user *UserInfo) (date time.Time, isErr bool) {

	if db == nil {
		Logs <- NewLog(nil, "SQL{GetFirstDate}", Error, sql_LostConnection)
		return time.Time{}, true
	}

	stmt := `
	with t as (select date as date 
		from operations
		where chat_id = $1 
		union 
		select date 
		from logs
		where chat_id = $1)

	select 
	DATE(min(date)) as date 
	from t;`

	err := db.Get(&date, stmt, user.ChatID)
	if err != nil {
		Logs <- NewLog(nil, "SQL{GetFirstDate}", Error, err.Error())
		return date, true
	}

	return date, false

}

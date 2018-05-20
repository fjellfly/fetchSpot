package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var connect = func(dbConfig store) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s:%d/%s", dbConfig.User, db.Password, db.Host, db.Port, db.DB))
	return
}

var getMessageKeys = func(db *sql.DB, unixTime int) (keys map[string]struct{}, err error) {

	rows, err := db.Query("SELECT id, messengerID FROM spot_messages")
	if err != nil {
		return
	}
	defer rows.Close()

	keys = make(map[string]struct{})

	for rows.Next() {
		var id int
		var messengerID string
		
		if err = rows.Scan(&id); err != nil {
			return
		}

		if err = rows.Scan(&messengerID); err != nil {
			return
		}

		keys[fmt.Sprintf("%d_%s", id, messengerID)] = struct{}{}
	}

	if rows.Err() != nil {
		return
	}

	return
}



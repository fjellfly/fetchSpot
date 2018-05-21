package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var prefix string

var connect = func(dbConfig store) (db *sql.DB, err error) {
	prefix = dbConfig.Prefix
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.TLS))
	return
}

// Get known messenger
var getMessenger = func(db *sql.DB) (messenger map[string]string, err error) {

	rows, err := db.Query(fmt.Sprintf("SELECT id, name FROM %s_messenger", prefix))
	if err != nil {
		return
	}
	defer rows.Close()

	messenger = make(map[string]string)

	for rows.Next() {
		var id string
		var name string

		if err = rows.Scan(&id, &name); err != nil {
			return
		}

		messenger[id] = name
	}

	err = rows.Err()
	return
}

// Get keys of all messenges not created before timeLimit
var getMessageKeys = func(db *sql.DB, timeLimit int64) (keys map[string]struct{}, err error) {

	rows, err := db.Query(fmt.Sprintf("SELECT id, messengerID FROM %s_messages WHERE unixTime >= ?", prefix), timeLimit)
	if err != nil {
		return
	}
	defer rows.Close()

	keys = make(map[string]struct{})

	for rows.Next() {
		var id int
		var messengerID string
		
		if err = rows.Scan(&id, &messengerID); err != nil {
			return
		}

		keys[fmt.Sprintf("%d_%s", id, messengerID)] = struct{}{}
	}

	err = rows.Err()
	return
}

var insertMessages = func(db *sql.DB, messages []message) (err error) {

	cmd, err := db.Prepare(fmt.Sprintf("INSERT INTO %s_messages SET id=?, messengerID=?, unixTime=?, messageType=?, location=POINT(?,?), dateTime=?, batteryState=?, messageContent=?, altitude=?", prefix))
	if err != nil {
		return
	}

	defer cmd.Close()

	for _, m := range messages {
		_, err = cmd.Exec(m.ID, m.MessengerID, m.UnixTime, m.MessageType, m.Longitude, m.Latitude, m.DateTime, m.BatteryState, m.MessageContent, m.Altitude)
		if err != nil {
			return
		}
	}

	return
}

var insertMessenger = func(db *sql.DB, messenger map[string]string) (err error) {

	cmd, err := db.Prepare(fmt.Sprintf("INSERT INTO %s_messenger SET id=?, name=?", prefix))
	if err != nil {
		return
	}
	defer cmd.Close()

	for id, name := range messenger {
		_, err = cmd.Exec(id, name)
		if err != nil {
			return
		}
	}

	return
}



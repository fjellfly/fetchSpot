package main

import (
	"fmt"
	"time"
	"database/sql"
)

// Get creation time of the oldest fetched message
func getTimeLimit(messages []message) (timeLimit int64) {
	timeLimit = time.Now().Unix()

	for _, m := range messages {
		if m.UnixTime < timeLimit {
			timeLimit = m.UnixTime
		}
	}

	return
}

// Return messages with are new and not already known to the database
var selectMessages = func(db *sql.DB, messages []message) (newMessages []message, err error) {

	timeLimit := getTimeLimit(messages)

	// Get keys (messageID+messengeID) of messages younger or equal aged then timeLimit from db
	keys, err := getMessageKeys(db, timeLimit)
	if err != nil {
		return
	}

	for _, message := range messages {

		if _, ok := keys[fmt.Sprintf("%d_%s", message.ID, message.MessengerID)]; ok {
			// Message with given key is already stored in db. skip.
			continue
		}

		newMessages = append(newMessages, message)
	}

	return
}

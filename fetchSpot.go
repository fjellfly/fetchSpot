package main

import (
	"fmt"
	"time"
)

type message struct {
	ID int64
	MessengerID string
	MessengerName string
	UnixTime int64
	MessageType string
	Latitude float64
	Longitude float64
	DateTime string
	BatteryState string
	MessageContent string
	Altitude float64
}

type  store struct {
	Host string
	Port int
	User string
	Name string
	Prefix string
	Password string
	TLS string
}

type feed struct {
	Password string
	ID string
}

func main(){

	// Import config
	feeds, dbConfig, err := readConfig("config.json")
	if err != nil {
		panic(fmt.Sprintf("Error while reading config file: %s", err.Error()))
	}

	// Fetch messages from spot
	var messages []message
	for _, feedInstance := range feeds {

		feedURL := getURL(feedInstance)

		feedMessages, err := getFeedMessages(feedURL)
		if err != nil {
			panic(fmt.Sprintf("Error while getting messages: %s", err.Error()))
		}

		messages = append(messages, feedMessages...)
	}

	// Quit if there are no new messages
	if len(messages) == 0 {
		return
	}

	// Time of creation of oldest message
	oldest := time.Now().Unix()
	for _, m := range messages {
		if m.UnixTime < oldest {
			oldest = m.UnixTime
		}
	}

	// Connect to db
	db, err := connect(dbConfig)
	if err != nil {
		panic(err.Error())
	}

	// Get keys (messageID+messengeID) of younger messages from db
	keys, err := getMessageKeys(db, 0)
	if err != nil {
		panic(err.Error())
	}

	// Select new messages
	newMessages := selectMessages(keys, messages)

	// Push new messages to db
	err = insertMessages(db, newMessages)
	if err != nil {
		panic(err.Error())
	}

	return
}

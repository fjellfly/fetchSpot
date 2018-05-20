package main

import (
	"fmt"
	"time"
)

type message struct {
	ID int
	MessengerID string
	MessengerName string
	UnixTime int
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
	DB string
	Prefix string
	Password string
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

	/*
	// Fetch messages from spot 
	var messages []message
	for _, feedInstance := range feeds {

		feedURL := getURL(feedInstance)
		fmt.Println(feedURL)
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
	oldest := time.Unix()
	for i, m := range messages {
		if m.UnixTime < oldest {
			oldest = m.UnixTime
		}
		fmt.Printf("#%d: %+v\n",i,m)
	}
	*/
	// Connect to db
	db, err := connect(dbConfig)
	if er != nil {
		panic(err.Error())
	}

	// Get keys (messageID+messengeID) of younger messages from db

	keys, err := getMessageKeys(db, 0)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%+v\n",keys)



	// Push new messages to database


	return
}

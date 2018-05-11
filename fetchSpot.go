package main

import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
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
	user string
	DB string
	Prefix string
	Password string
}

type feed struct {
	Password string
	ID string
}

type feedBasicResponse struct {
	Count int
	TotalCount int
	ActivityCount int
}

var readConfig = func(file string) (feeds []feed, mysql store, err error) {

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	config := &struct {
		Mysql store
		Feeds []feed
	}{}

	err = json.Unmarshal(fileContent, config)
	if err != nil {
		return
	}

	feeds = config.Feeds
	mysql = config.Mysql

	return
}

var getURL = func(feedInstance feed) string {

	u := url.URL{
		Scheme: "https",
		Host: "api.findmespot.com",
		Path: fmt.Sprintf("spot-main-web/consumer/rest-api/2.0/public/feed/%s/message.json",feedInstance.ID),
	}

	if feedInstance.Password != "" {
		u.RawQuery = fmt.Sprintf("feedPassword=%s",feedInstance.Password)
	}

	return u.String()
}

var doGetFeed = func(feedURL string) (*http.Response, error) {
	return http.Get(feedURL)
}

var getMessages = func(feedURL string) (messages []message, err error) {

	httpResponse, err := doGetFeed(feedURL)
	if err != nil {
		return
	}

	defer httpResponse.Body.Close()
	rawResponse, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	// Determine number of messages
	basicResponse := &struct{
		Response struct {
			FeedMessageResponse feedBasicResponse
			Errors *struct{
				Error struct {
					Code string
					Text string
					Description string
				}
			}
		}
	}{}

	err = json.Unmarshal(rawResponse, &basicResponse)
	if err != nil {
		return
	}

	// Check for errors within the response
	if basicResponse.Response.Errors != nil {

		switch basicResponse.Response.Errors.Error.Code {
		case "E-0195":
			// This is not an error, just a "nothing to do".
			return
		default:
			// Return an error
			err = fmt.Errorf("%s (%s): %s",basicResponse.Response.Errors.Error.Text,
				basicResponse.Response.Errors.Error.Code, basicResponse.Response.Errors.Error.Description)
			return
		}
	}

	log.Printf("%d messages found\n",basicResponse.Response.FeedMessageResponse.Count)

	// Get messages
	switch basicResponse.Response.FeedMessageResponse.Count {
	case 1:
		// Fetch single message
		singleMessageResponse := &struct{
			Response struct {
				FeedMessageResponse struct{
					Messages struct {
						Message message
					}
				}
			}
		}{}

		err = json.Unmarshal(rawResponse, singleMessageResponse)
		if err !=nil {
			err = fmt.Errorf("Error while decoding a single message: %s", err.Error())
			return
		}

		messages = append(messages,singleMessageResponse.Response.FeedMessageResponse.Messages.Message)

	default:
		// Fetch multiple messages
		multipleMessageResponse := &struct{
			Response struct {
				FeedMessageResponse struct{
					Messages struct {
						Message []message
					}
				}
			}
		}{}

		err = json.Unmarshal(rawResponse, multipleMessageResponse)
		if err !=nil {
			err = fmt.Errorf("Error while decoding multiple messages: %S", err.Error())
		}

		messages = multipleMessageResponse.Response.FeedMessageResponse.Messages.Message
	}

	return
}

func main(){

	feeds, mysql, err := readConfig("config.json")
	if err != nil {
		panic(fmt.Sprintf("Error while reading config file: %s", err.Error()))
	}

	var messages []message

	for _, feedInstance := range feeds {

		feedURL := getURL(feedInstance)
		fmt.Println(feedURL)
		feedMessages, err := getMessages(feedURL)
		if err != nil {
			panic(fmt.Sprintf("Error while getting messages: %s", err.Error()))
		}

		messages = append(messages, feedMessages...)
	}

	if len(messages) == 0 {
		log.Println("No messages")
	}

	for i, m := range messages {
		fmt.Printf("#%d: %+v\n",i,m)
	}

	return
}

package main

import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"flag"
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

type feedBasicResponse struct {
	Count int
	TotalCount int
	ActivityCount int
}

func main(){

	password := flag.String("password", "", "password for selected feed (optional)")
	feedID := flag.String("feed", "", "ID of feed to fetch messages from")

	flag.Parse()

	if *feedID == "" {
		panic("no feed given")
	}

	if *password != "" {
		*password = fmt.Sprintf("feedPassword=%s",*password)
	}

	u := url.URL{
		Scheme: "https",
		Host: "api.findmespot.com",
		Path: fmt.Sprintf("spot-main-web/consumer/rest-api/2.0/public/feed/%s/message.json",*feedID),
		RawQuery: *password,
	}

	responseReader, err := http.Get(u.String())
	if err != nil {
		panic (err.Error())
	}

	defer responseReader.Body.Close()
	rawResponse, err := ioutil.ReadAll(responseReader.Body)
	if err != nil {
		panic (err.Error())
	}

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
		panic(err.Error())
	}

	// Check for errors within the response
	if basicResponse.Response.Errors != nil {

		switch basicResponse.Response.Errors.Error.Code {
		case "E-0195":
			log.Println("No messages availible")
			return
		default:
			panic(fmt.Sprintf("%s (%s): %s",basicResponse.Response.Errors.Error.Text, 
				basicResponse.Response.Errors.Error.Code, basicResponse.Response.Errors.Error.Description))
		}
	}

	log.Printf("%d messages found\n",basicResponse.Response.FeedMessageResponse.Count)

	var messages []message

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

		err := json.Unmarshal(rawResponse, singleMessageResponse)
		if err !=nil {
			panic("Error while decoding a single message")
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

		err := json.Unmarshal(rawResponse, multipleMessageResponse)
		if err !=nil {
			panic("Error while decoding multiple messages")
		}

		messages = multipleMessageResponse.Response.FeedMessageResponse.Messages.Message
	}

	for i, m := range messages {
		fmt.Printf("#%d: %+v\n",i,m)
	}

	return
}

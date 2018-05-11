package main

import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type feedBasicResponse struct {
	Count int
	TotalCount int
	ActivityCount int
}

type responseError struct {
	Code string
	Text string
	Description string
}

func (e *responseError) Error() string {
	return fmt.Sprintf("%s (%s): %s", e.Text, e.Code, e.Description)
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

var doGetFeed = func(feedURL string) (rawResponse []byte, err error) {

	httpResponse, err := http.Get(feedURL)
	if err != nil {
		return
	}

	defer httpResponse.Body.Close()

	rawResponse, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	return
}

var getFeedMessages = func(feedURL string) (messages []message, err error) {

	rawResponse, err := doGetFeed(feedURL)
	if err != nil {
		return
	}

	// Determine number of messages
	basicResponse := &struct{
		Response struct {
			FeedMessageResponse feedBasicResponse
			Errors *struct{
				Error responseError
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
			err = &responseError{
				Code: basicResponse.Response.Errors.Error.Code,
				Description: basicResponse.Response.Errors.Error.Description,
				Text: basicResponse.Response.Errors.Error.Text,
			}

			return
		}
	}

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
			err = fmt.Errorf("Error while decoding multiple messages: %s", err.Error())
		}

		messages = multipleMessageResponse.Response.FeedMessageResponse.Messages.Message
	}

	return
}

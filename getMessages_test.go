package main

import (
	"testing"
	"fmt"
	"sort"
)

func TestGetUrl(t *testing.T) {

	expectedURLFront := "https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed"

	testCases := []struct{
		desc string
		expectedURL string
		feedInstance feed
	}{
		{
			desc: "Feed with password",
			expectedURL: fmt.Sprintf("%s/%s/message.json?feedPassword=%s",expectedURLFront, "testFeedWithPassword", "testPassword"),
			feedInstance: feed{
				ID: "testFeedWithPassword",
				Password: "testPassword",
			},
		},
		{
			desc: "Feed without password",
			expectedURL: fmt.Sprintf("%s/%s/message.json",expectedURLFront, "testFeed"),
			feedInstance: feed{
				ID: "testFeed",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			url := getURL(tc.feedInstance)

			if url != tc.expectedURL {
				t.Errorf("Expected %s but got %s", tc.expectedURL, url)
			}
		})
	}

	return
}

func TestGetFeedMessages(t *testing.T) {

	testCases := []struct{
		desc string
		payload []byte
		expectedResponseError *responseError
		expectedMessages []message
	}{
		{
			desc: "No new messages",
			payload: []byte(`{
				"Response": {
					"FeedMessageResponse": {
						"Count": 0,
						"TotalCount": 42,
						"ActivityCount": 0
					},
					"Errors": {
						"Error": {
							"Code": "E-0195"
						}
					}
				}
			}`),
		},
		{
			desc: "Error while fetching messages",
			payload: []byte(`{
				"Response": {
					"FeedMessageResponse": {
						"Count": 20
					},
					"Errors": {
						"Error": {
							"Code": "E-6666"
						}
					}
				}
			}`),
			expectedResponseError: &responseError{
					Code: "E-6666",
			},
		},
		{
			desc: "One new message",
			payload: []byte(`{
				"Response": {
					"FeedMessageResponse": {
						"Count": 1,
						"Messages": {
							"Message": {
								"ID": 42,
								"MessageContent": "testMessage"
							}
						}
					}
				}
			}`),
			expectedMessages: []message{
				{
					ID: 42,
					MessageContent: "testMessage",
				},
			},
		},
		{
			desc: "Two (multiple) new messages",
			payload: []byte(`{
				"Response": {
					"FeedMessageResponse": {
						"Count": 2,
						"Messages": {
							"Message": [
								{
									"ID": 42,
									"MessageContent": "testMessage"
								},
								{
									"ID": 43,
									"MessageContent": "anotherTestMessage"
								}
							]
						}
					}
				}
			}`),
			expectedMessages: []message{
				{
					ID: 42,
					MessageContent: "testMessage",
				},
				{
					ID: 43,
					MessageContent: "anotherTestMessage",
				},
			},
		},
	}

	doGetFeedBackup := doGetFeed
	defer func(){
		doGetFeed = doGetFeedBackup
	}()

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			doGetFeed = func(_ string) ([]byte, error) {
				return tc.payload, nil
			}

			messages, err := getFeedMessages("testURL")
			if err != nil {
				if respErr, ok := err.(*responseError); ok {
					if tc.expectedResponseError != nil {
						if respErr.Code != tc.expectedResponseError.Code {
							t.Fatalf("Expected Error-Code to be %s but got %s", tc.expectedResponseError.Code, respErr.Code)
						}

						// Expected response error - erverything is fine, return
						return
					}
					t.Fatalf("Got unexpected response-error: %s", err.Error())
				}

				t.Fatalf("Got unexpected error: %s", err.Error())
			}

			if len(messages) != len(tc.expectedMessages) {
				t.Errorf("Expected %d message(s) but got %d", len(tc.expectedMessages), len(messages))
			}

			sort.Slice(messages, func(i,j int) bool { return messages[i].ID < messages[j].ID })
			sort.Slice(tc.expectedMessages,  func(i,j int) bool {return tc.expectedMessages[i].ID < tc.expectedMessages[j].ID })

			for i := range messages {
				if messages[i].ID != tc.expectedMessages[i].ID || messages[i].MessageContent != tc.expectedMessages[i].MessageContent {
					t.Errorf(fmt.Sprintf("Unexpected message: Expected %+v got %+v", tc.expectedMessages[i], messages[i]))
				}
			}

		});
	}

	return
}

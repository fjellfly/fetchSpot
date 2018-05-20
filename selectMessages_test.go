package main

import (
	"testing"
	"fmt"
	"database/sql"
)

func TestSelectMessages(t *testing.T) {

	testCases := []struct {
		desc string
		expectedTimeLimit int64
		fetchedMessages []message
		keysOfOldMessages map[string]struct{}
		keysOfNewMessages map[string]struct{}
	}{
		{
			desc: "New messages only",
			fetchedMessages: []message{
				{
					ID: 12,
					UnixTime: 900,
					MessengerID: "TestMessenger1",
				},
				{
					ID: 13,
					UnixTime: 1000,
					MessengerID: "TestMessenger1",
				},
				{
					ID: 13,
					UnixTime: 1000,
					MessengerID: "TestMessenger2",
				},
			},
			keysOfNewMessages: map[string]struct{}{
				"12_TestMessenger1": struct{}{},
				"13_TestMessenger1": struct{}{},
				"13_TestMessenger2": struct{}{},
			},
			expectedTimeLimit: 900,
		},
		{
			desc: "Old messages only",
			fetchedMessages: []message{
				{
					ID: 12,
					UnixTime: 900,
					MessengerID: "TestMessenger1",
				},
				{
					ID: 13,
					UnixTime: 1000,
					MessengerID: "TestMessenger1",
				},
				{
					ID: 13,
					UnixTime: 1000,
					MessengerID: "TestMessenger2",
				},
			},
			keysOfOldMessages: map[string]struct{}{
				"12_TestMessenger1": struct{}{},
				"13_TestMessenger1": struct{}{},
				"13_TestMessenger2": struct{}{},
			},
			expectedTimeLimit: 900,
		},
		{
			desc: "Old and new messages",
			fetchedMessages: []message{
				{
					ID: 12,
					UnixTime: 900,
					MessengerID: "TestMessenger1",
				},
				{
					ID: 13,
					UnixTime: 1000,
					MessengerID: "TestMessenger1",
				},
				{
					ID: 13,
					UnixTime: 1000,
					MessengerID: "TestMessenger2",
				},
			},
			keysOfOldMessages: map[string]struct{}{
				"12_TestMessenger1": struct{}{},
			},
			keysOfNewMessages: map[string]struct{}{
				"13_TestMessenger1": struct{}{},
				"13_TestMessenger2": struct{}{},
			},
			expectedTimeLimit: 900,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T){

			getMessageKeysBackup := getMessageKeys
			defer func() {
				getMessageKeys = getMessageKeysBackup
			}()

			getMessageKeys = func(db *sql.DB, timeLimit int64) (keys map[string]struct{}, err error) {
				if timeLimit != tc.expectedTimeLimit {
					t.Errorf("Expected timeLimit to be %d but got %d", tc.expectedTimeLimit, timeLimit)
				}

				keys = tc.keysOfOldMessages
				return
			}

			db := &sql.DB{}

			newMessages, err := selectMessages(db, tc.fetchedMessages)
			if err != nil {
				t.Errorf("Got unexpected error: %s", err.Error())
			}

			for _, m := range newMessages {

				key := fmt.Sprintf("%d_%s", m.ID, m.MessengerID)
				if _, ok := tc.keysOfNewMessages[key]; !ok {
					t.Errorf("Got unexpected message with ID %s", key)
				}

				delete(tc.keysOfNewMessages, key)
			}

			if len(tc.keysOfNewMessages) != 0 {
				var keys []string

				for key, _ := range tc.keysOfNewMessages {
					keys = append(keys, key)
				}

				t.Errorf("Didn't get a message for the following key(s): %s)", keys)
			}
		})
	}

	return
}

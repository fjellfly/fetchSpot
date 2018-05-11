package main

import (
	"fmt"
)

var selectMessages = func(keys map[string]struct{}, messages []message) (newMessages []message) {

	for _, message := range messages {

		if _, ok := keys[fmt.Sprintf("%d_%s", message.ID, message.MessengerID)]; ok {
			// Message with given key is already stored in db. skip.
			continue
		}

		newMessages = append(newMessages, message)
	}

	return
}

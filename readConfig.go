package main

import (
	"io/ioutil"
	"encoding/json"
)

var doReadConfig = func(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

var readConfig = func(file string) (feeds []feed, mysql store, err error) {

	fileContent, err := doReadConfig(file)
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

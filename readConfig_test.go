package main

import(
	"testing"
)

func TestReadConfig(t *testing.T) {

	doReadConfigBackup := doReadConfig
	defer func(){
		doReadConfig = doReadConfigBackup
	}()

	doReadConfig = func(_ string)(content []byte, err error){
		content = []byte(`{ 
			"mysql": { 
				"port": 42, 
				"prefix": "spot"
			},
			"feeds": [
				{
					"password": "passwordA",
					"id": "feedA"
				},
				{
					"password": "passwordB",
					"id": "feedB"
				}
			]
		}`)

		return
	}

	feeds, mysql, err := readConfig("test")

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if len(feeds)!=2 {
		t.Errorf("Expected two feeds but got %d", len(feeds))
	}

	if mysql.Port != 42 {
		t.Errorf("Expected port to be '42' but got %d", mysql.Port)
	}

	if mysql.Prefix != "spot" {
		t.Errorf("Expected prefix to be 'spot' but got %s", mysql.Prefix)
	}

	return
}




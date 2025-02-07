package testutils

import (
	"encoding/json"
	"log"
)

func ErrorStringFromBody(body []byte) string {
	var errJson map[string]string
	err := json.Unmarshal(body, &errJson)
	if err != nil {
		log.Fatal(err) // it's safe since it's used only in tests
	}
	return errJson["error"]
}

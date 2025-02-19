package testutils

import (
	"encoding/json"
	"log"
)

func ErrorStringFromBody(body []byte) string {
	var errJson map[string]any
	err := json.Unmarshal(body, &errJson)
	if err != nil {
		// it's safe since it's used only in tests
		log.Fatal("ERROR [testutils]: ErrorStringFromBody json.Unmarshal err=", err)
	}
	errVal := errJson["error"]
	switch errVal := errVal.(type) {
	case string:
		return errVal
	default:
		return ""
	}
}

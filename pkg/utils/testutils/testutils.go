package testutils

import (
	"encoding/json"
	"log"
)

func ErrorCodeFromBody(body []byte) string {
	var resp map[string]any
	err := json.Unmarshal(body, &resp)
	if err != nil {
		// it's safe since it's used only in tests
		log.Fatal("ERROR [testutils]: ErrorStringFromBody json.Unmarshal err=", err)
	}
	errCode := resp["error_code"]
	switch code := errCode.(type) {
	case string:
		return code
	default:
		return ""
	}
}

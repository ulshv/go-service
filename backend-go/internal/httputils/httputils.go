package httputils

import (
	"encoding/json"
	"net/http"
)

var (
	errJsonDecode = "error while decoding json, make sure the data structure is correct"
)

func ErrorJson(w http.ResponseWriter, msg string, statusCode int) {
	errJson, _ := json.Marshal(map[string]string{"error": msg})
	w.WriteHeader(statusCode)
	w.Write(errJson)
}

func DecodeJsonAndHandleErr(w http.ResponseWriter, r *http.Request, dto any) error {
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		ErrorJson(w, errJsonDecode, http.StatusBadRequest)
		return err
	}
	return nil
}

func WriteJson(w http.ResponseWriter, data any) {
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

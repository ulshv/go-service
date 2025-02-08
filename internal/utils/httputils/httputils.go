package httputils

import (
	"encoding/json"
	"net/http"

	logger_mod "github.com/ulshv/go-service/internal/logger"
)

var (
	errJsonDecode = "error while decoding json, make sure the data structure is correct"
	logger        = logger_mod.NewLogger("httputils")
)

func WriteJson(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		msg := "error while encoding json response"
		logger.Error(msg, "error", err)
		WriteErrorJson(w, msg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func WriteErrorJson(w http.ResponseWriter, msg string, statusCode int) {
	errJson, _ := json.Marshal(map[string]string{"error": msg})
	w.WriteHeader(statusCode)
	w.Write(errJson)
}

func DecodeBody(w http.ResponseWriter, r *http.Request, dto any) error {
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		WriteErrorJson(w, errJsonDecode, http.StatusBadRequest)
		return err
	}
	return nil
}

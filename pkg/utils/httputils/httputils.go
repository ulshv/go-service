package httputils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ulshv/go-service/internal/core/httperrs"
	"github.com/ulshv/go-service/pkg/logs"
)

var (
	errJsonEncode = errors.New("error while encoding json response")
	errJsonDecode = errors.New("error while decoding json, make sure the data structure is correct")
	logger        = logs.NewLogger("httputils")
)

func WriteJson(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error(errJsonEncode.Error(), "error", err)
		WriteErrorJson(w, errJsonEncode, httperrs.ErrCodeInternal, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func WriteErrorJson(w http.ResponseWriter, err error, errCode string, statusCode int) {
	// TODO - check for err's interface, if it's pq/sqlite SQL error - don't return the text on production, only in dev mode
	errJson, _ := json.Marshal(map[string]string{"error": err.Error(), "error_code": errCode})
	w.WriteHeader(statusCode)
	w.Write(errJson)
}

func DecodeBody(w http.ResponseWriter, r *http.Request, dto any) error {
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		WriteErrorJson(w, errJsonDecode, httperrs.ErrCodeMalformedJson, http.StatusBadRequest)
		return err
	}
	return nil
}
